/*  doggleship.go--salvo between programs
    Written 2022 by Eric Olson

    Version 1 run SALVO from 101 Basic Computer Games
    Version 2 cleanup running programs on interrupt
    Version 3 add Fortran driver and refactor 
    Version 4 add Ebiten graphical dashboard
    Version 5 add option to turn off dashboard
    Version 6 add mutex to avoid racing with monitor
    Version 7 don't forget to reap the zombies */

package main

import (
    "C"; "bufio"; "flag"; "fmt"; "os"; "os/exec"; "os/signal"
    "reflect"; "runtime"; "strconv"; "strings"; "sync"
    "syscall"; "time"; "unsafe"
)

type (
    stspec int
    cospec struct {
        i,j int
    }
    btspec struct {
        s string
        p []cospec
    }
    plfunc func(b []btspec)
    ogfunc func(r []cospec)
    icfunc func()([]cospec,stspec)
    bofunc func(game *gmspec,cmd,log string)
    bospec struct {
        boot bofunc
        pname string
    }
    srspec struct {
        s []string
        mu sync.Mutex
    }
    gmspec struct {
        b []btspec
        sr srspec
        pr *os.Process
        fd *os.File
        first bool
        cmd,log string
        place plfunc
        incoming icfunc
        outgoing ogfunc
    }    
)

const (
    running stspec=iota; lost; won
    twait=time.Duration(1)
    tmult=1
)

var (
    game1,game2 gmspec
    trial,turn int
    win1,win2 int
    dover,tnum int
    batch bool
)

func ishere(g *gmspec,i,j int)bool {
    i++; j++
    for k:=0;k<len(g.b);k++ {
        for l:=0;l<len(g.b[k].p);l++ {
            if g.b[k].p[l].i==i&&g.b[k].p[l].j==j {
                return true
            }
        }
    }
    return false
}

func cleanup(){
    if game1.pr!=nil {
        game1.pr.Signal(syscall.SIGINT)
        game1.pr=nil
    }
    if game2.pr!=nil {
        game2.pr.Signal(syscall.SIGINT)
        game2.pr=nil
    }
}

func onexit(){
    c:=make(chan os.Signal,1)
    signal.Notify(c,syscall.SIGINT,syscall.SIGHUP,syscall.SIGQUIT)
    s:=<-c
    fmt.Printf("Exiting on signal %d...\n",s)
    cleanup()
    os.Exit(1)
}

func myexit(r int){
    cleanup()
    os.Exit(r)
}

func spawn(cl string)(*os.File,*os.Process) {
    fmaster,err:=os.OpenFile("/dev/ptmx",os.O_RDWR,0)
    if err!=nil {
        fmt.Printf("Error opening /dev/ptms for read and write!\n")
        myexit(1)
    }
    var pn string
    {
        var n C.uint
        syscall.Syscall(syscall.SYS_IOCTL,fmaster.Fd(),
            syscall.TIOCGPTN,uintptr(unsafe.Pointer(&n)))
        pn=fmt.Sprintf("/dev/pts/%d",n)
        n=0
        syscall.Syscall(syscall.SYS_IOCTL,fmaster.Fd(),
            syscall.TIOCSPTLCK,uintptr(unsafe.Pointer(&n)))
    }
    fslave,err:=os.OpenFile(pn,os.O_RDWR,0)
    if err!=nil {
        fmt.Printf("Error opening %s for read and write\n%v!\n",pn,err)
        myexit(1)
    }
    cmd:=exec.Command(cl)
    cmd.Stdin=fslave; cmd.Stdout=fslave; cmd.Stderr=fslave
    cmd.Start()
    go cmd.Wait()
    return fmaster,cmd.Process
}

func (sr *srspec)mulen()int{
    sr.mu.Lock()
    q:=len(sr.s)
    sr.mu.Unlock()
    return q
}

func (sr *srspec)mustr(p int)string{
    sr.mu.Lock()
    s:=""
    if p<len(sr.s) {
        s=sr.s[p]
    } 
    sr.mu.Unlock()
    return s
}

func monitor(sr *srspec,fp *os.File,log string) {
    logio,err:=os.Create(log)
    if err!=nil {
        fmt.Printf("Error opening %s for write!\n",log)
        myexit(1)
    }
    defer logio.Close()
    rd:=bufio.NewReader(fp)
    for {
        var b [128]byte
        n,err:=rd.Read(b[:])
        if err!=nil { return }
        var i,j int
        for i,j=0,0;j<n;j++ {
            if b[j]=='\n' {
                i=j+1
            } else if b[j]=='\r' {
                fmt.Fprintf(logio,"%s\n",b[i:j])
                sr.mu.Lock()
                sr.s[len(sr.s)-1]+=string(b[i:j])
                sr.s=append(sr.s,"")
                sr.mu.Unlock()
                i=j+1
            }
        }
        if i<j {
            fmt.Fprintf(logio,"%s",b[i:j])
            sr.mu.Lock()
            sr.s[len(sr.s)-1]+=string(b[i:j])
            sr.mu.Unlock()
        }
    }
}

func istail(a,b string)bool {
    if len(a)<len(b) { return false }
    if a[len(a)-len(b):]==b { return true }
    return false
}

func ishead(a,b string)bool {
    if len(a)<len(b) { return false }
    if a[0:len(b)]==b { return true }
    return false
}

func (game *gmspec)mkgame(first bool) {
    if game.b==nil {
        game.b=[]btspec{
            btspec{"BATTLESHIP",make([]cospec,0,5)},
            btspec{"CRUISER",make([]cospec,0,5)},
            btspec{"DESTROYER<A>",make([]cospec,0,5)},
            btspec{"DESTROYER<B>",make([]cospec,0,5)},
        }
    } else {
        for i:=0;i<len(game.b);i++ {
            game.b[i].p=game.b[i].p[:0]
        }
    }
    if game.sr.s==nil {
        game.sr.s=make([]string,1,512)
    } else {
        game.sr.s=game.sr.s[:1]
        game.sr.s[0]=""
    }    
    game.first=first
}

func prboats(b []btspec){
    fmt.Printf("\n")
    for i:=0;i<len(b);i++ {
        fmt.Printf("|%s\n",b[i].s)
        for j:=0;j<len(b[i].p);j++ {
            fmt.Printf("| %d  %d\n",b[i].p[j].i,b[i].p[j].j)
        }
    }
}

var winlose=[]string{
    "I HAVE WON","YOU HAVE WON"," YOU WIN!"," YOU LOSE!",
}

func isend(s string)bool {
    for i:=0;i<len(winlose);i++ {
        if ishead(s,winlose[i]) { return true }
    }
    return false
}

func (game *gmspec)mufindtail(s string)bool {
    found:=false
    dotail:=func()bool {
        p:=len(game.sr.s)-1
        if p>=0 {
            if istail(game.sr.s[p],s){ found=true; return true }
            if game.sr.s[p]==">" { found=false; return true }
        }
        if p>2 {
            for q:=p-2;q<=p;q++ {
                if isend(game.sr.s[q]) {
                    found=false
                    return true
                }
            }
        }
        return false
    }
    t:=twait
    for {
        game.sr.mu.Lock()
        r:=dotail()
        game.sr.mu.Unlock()
        if r { return found }
        time.Sleep(t*time.Millisecond)
        if t<128 { t*=2 }
    }
}

func (game *gmspec)munextprompt(p int,s string)int {
    t:=twait
    for {
        game.mufindtail(s)
        r:=game.sr.mulen()
        if r>p { return r }
        time.Sleep(t*time.Millisecond)
        if t<128 { t*=2 }
    }
}

func (game *gmspec)muyouhave()int{
    match:=func()int {
        p:=len(game.sr.s)
        for q:=p-3;q<p;q++ {
            if isend(game.sr.s[q]) { return q }
            if strings.Contains(game.sr.s[q],"AND YOU HAVE") { return q }
        }
        return 0
    }
    t:=twait
    for {
        game.sr.mu.Lock()
        p:=match()
        game.sr.mu.Unlock()
        if p>0 { return p }
        time.Sleep(t*time.Millisecond)
        if t<128 { t*=2 }
    }
}

func bbcboot(game *gmspec,cmd,log string) {
    getname:=func(s string)int {
        for j:=0;j<len(game.b);j++ {
            if s==game.b[j].s {
                return j
            }
        }
        return -1
    }
    game.cmd=cmd; game.log=log
    game.place=game.bbcplace
    game.incoming=game.bbcincoming
    game.outgoing=game.bbcoutgoing
    game.fd,game.pr=spawn(cmd)
    go monitor(&game.sr,game.fd,log)
    game.mufindtail("DO YOU WANT TO START? ")
    p:=game.sr.mulen()
    fmt.Fprintf(game.fd,"WHERE ARE YOUR SHIPS?\n")
    game.munextprompt(p,"? ")
    j:=-1
    game.sr.mu.Lock()
    for i:=p;i<len(game.sr.s);i++ {
        r:=getname(game.sr.s[i])
        if r>=0 { j=r
        } else if j>=0{
            sxy:=strings.Fields(game.sr.s[i])
            if len(sxy)==2 {
                x,err:=strconv.Atoi(sxy[0])
                if err!=nil { continue }
                y,err:=strconv.Atoi(sxy[1])
                if err!=nil { continue }
                game.b[j].p=append(game.b[j].p,cospec{x,y})
            }
        }
    }
    game.sr.mu.Unlock()
    if game.first {
        fmt.Fprintf(game.fd,"NO\n")
    } else {
        fmt.Fprintf(game.fd,"YES\n")
    }
}

func forboot(game *gmspec,cmd,log string) {
    game.cmd=cmd; game.log=log
    game.place=game.forplace
    game.incoming=game.forincoming
    game.outgoing=game.foroutgoing
    game.fd,game.pr=spawn(cmd)
    go monitor(&game.sr,game.fd,log)
    game.mufindtail("CHOICE: ")
    p:=len(game.sr.s)
    fmt.Fprintf(game.fd,"0\n")
    game.munextprompt(p,": ")
    var (them [10][10]int; q=0)
    game.sr.mu.Lock()
    for i:=p;i<len(game.sr.s);i++ {
        if q>0{
            if istail(game.sr.s[i],"-----"){ break }
            xs:=strings.Fields(game.sr.s[i])
            for j:=0;j<len(xs);j++ {
                them[i-q][j],_=strconv.Atoi(xs[j])
            }
        }
        if istail(game.sr.s[i],"-----"){ 
            q=i+1
        }
    }
    game.sr.mu.Unlock()
    for i:=0;i<10;i++ {
        for j:=0;j<10;j++ {
            q=-them[i][j]-1
            if q>=0 {
                game.b[q].p=append(game.b[q].p,cospec{j+1,i+1})
            }
        }
    }
    if game.first {
        fmt.Fprintf(game.fd,"2\n")
    } else {
        fmt.Fprintf(game.fd,"1\n")
    }
}

func (game *gmspec)bbcplace(b []btspec){
    p:=game.sr.mulen()-2
    for i:=0;i<len(b);i++ {
        for j:=0;j<len(b[i].p);j++ {
            p=game.munextprompt(p,"? ")
            fmt.Fprintf(game.fd,"%d,%d\n",
                b[i].p[j].i,b[i].p[j].j)
        }
    }
    game.mufindtail("DO YOU WANT TO SEE MY SHOTS? ")
    p=game.sr.mulen()
    fmt.Fprintf(game.fd,"YES\n")
    game.munextprompt(p,"? ")
}

var dirs=[8]cospec{
    cospec{0,1},cospec{-1,1},cospec{-1,0},cospec{-1,-1},
    cospec{0,-1},cospec{1,-1},cospec{1,0},cospec{1,1}}

func (game *gmspec)forplace(b []btspec){
    p:=game.sr.mulen()-2
    for i:=0;i<len(b);i++ {
        p=game.munextprompt(p,": ")
        fmt.Fprintf(game.fd,"%d%d\n",
            b[i].p[0].i-1,b[i].p[0].j-1)
        p=game.munextprompt(p,": ")
        d:=cospec{b[i].p[1].i-b[i].p[0].i,b[i].p[1].j-b[i].p[0].j}
        var j int
        for j=0;j<len(dirs);j++ {
            if dirs[j]==d { break }
        }
        fmt.Fprintf(game.fd,"%d\n",j+1)
    }
    p=game.muyouhave()
    game.munextprompt(p,": ")
}

func (game *gmspec)bbcoutgoing(r []cospec){
    p:=game.sr.mulen()-2
    for i:=0;i<len(r);i++ {
        p=game.munextprompt(p,"? ")
        fmt.Fprintf(game.fd,"%d,%d\n",r[i].i,r[i].j)
    }
    game.munextprompt(p,"? ")
}

func (game *gmspec)foroutgoing(r []cospec){
    p:=game.sr.mulen()-2
    for i:=0;i<len(r);i++ {
        p=game.munextprompt(p,": ")
        fmt.Fprintf(game.fd,"%d%d\n",r[i].i-1,r[i].j-1)
    }
    game.muyouhave()
    game.munextprompt(p,": ")
}

func (game *gmspec)bbcincoming()([]cospec,stspec) {
    r:=make([]cospec,0,7)
    var (p int; sn []string)
    getst:=func()stspec {
        for p=len(game.sr.s)-1;p>=0;p-- {
            if ishead(game.sr.s[p],"I HAVE WON"){
                return won
            } else if ishead(game.sr.s[p],"YOU HAVE WON"){
                return lost
            } else if ishead(game.sr.s[p],"I HAVE "){
                sn=strings.Fields(game.sr.s[p])
                break
            }
        }
        return running
    }
    game.sr.mu.Lock()
    st:=getst()
    game.sr.mu.Unlock()
    if st!=running { return r,st }
    if len(sn)!=4 {
        fmt.Printf("%s:Couldn't parse number of shots from '%s'\n",
            game.log,game.sr.s[p])
        myexit(1)
    }
    n,err:=strconv.Atoi(sn[2])
    if err!=nil {
        fmt.Printf("%s:Couldn't convert '%s' to number of shots\n",
            game.log,sn[2])
        myexit(1)
    }
    p++; n+=p
    for i:=p;i<n;i++ {
        sxy:=strings.Fields(game.sr.mustr(i))
        if len(sxy)==2 {
            x,err:=strconv.Atoi(sxy[0])
            if err!=nil { continue }
            y,err:=strconv.Atoi(sxy[1])
            if err!=nil { continue }
            r=append(r,cospec{x,y})
        }
    }
    return r,running
}

func (game *gmspec)forincoming()([]cospec,stspec) {
    r:=make([]cospec,0,7)
    var (p int; sn []string)
    getst:=func()stspec {
        for p=len(game.sr.s)-1;p>=0;p-- {
            if ishead(game.sr.s[p]," YOU LOSE!"){
                return won
            } else if ishead(game.sr.s[p]," YOU WIN!"){
                return lost
            } else if strings.Contains(game.sr.s[p],"AND I HAVE"){
                sn=strings.Fields(game.sr.s[p])
                break
            }
        }
        return running
    }
    game.sr.mu.Lock()
    st:=getst()
    game.sr.mu.Unlock()
    if st!=running { return r,st }
    if len(sn)!=9 {
        fmt.Printf("%s:Couldn't parse number of shots from '%s'\n",
            game.log,game.sr.s[p])
        myexit(1)
    }
    n,err:=strconv.Atoi(sn[7])
    if err!=nil {
        fmt.Printf("%s:Couldn't convert '%s' to number of shots\n",
            game.log,sn[2])
        myexit(1)
    }
    p+=2; n+=p
    for i:=p;i<n;i++ {
        sxy:=strings.Fields(game.sr.mustr(i))
        if len(sxy)==3 {
            xy,err:=strconv.Atoi(sxy[2])
            if err!=nil { continue }
            r=append(r,cospec{xy/10+1,xy%10+1})
        }
    }
    return r,running
}

func doerase() {
    dover=0
    turn=0
    for i:=0;i<10;i++ {
        for j:=0;j<10;j++ {
            bdleft[i][j]=0
        }
    }
    for i:=0;i<10;i++ {
        for j:=0;j<10;j++ {
            bdright[i][j]=0
        }
    }
}

func dotrial(bt1,bt2 bospec){
    doerase()
    game1.mkgame(true)
    game2.mkgame(false)
    bt1.boot(&game1,bt1.pname,fmt.Sprintf("s%04d_1.txt",trial))
    bt2.boot(&game2,bt2.pname,fmt.Sprintf("s%04d_2.txt",trial))
    game1.place(game2.b)
    game2.place(game1.b)
    dorefresh(tmult*2*time.Second)
    for turn=1;;turn++ {
        r1,status:=game1.incoming()
        if status==running {
            for k:=0;k<len(r1);k++ {
                i:=r1[k].i; j:=r1[k].j
                bdleft[i-1][j-1]=turn
            }
            game2.outgoing(r1)
        } else {
            if status==lost {
                dover=2; win2++
            } else {
                dover=1; win1++
            }
            break
        } 
        r2,status:=game2.incoming()
        if status==running {
            for k:=0;k<len(r2);k++ {
                i:=r2[k].i; j:=r2[k].j
                bdright[i-1][j-1]=turn
            }
            game1.outgoing(r2)
        } else {
            if status==lost {
                dover=1; win1++
            } else {
                dover=2; win2++
            }
            break
        } 
        dorefresh(tmult*time.Second)
//        fmt.Printf("Press Enter: ")
//        var input string
//        fmt.Scanln(&input)
    }
    fmt.Printf("%7d %7d %7d\n",trial,dover,turn)
    dorefresh(tmult*4*time.Second)
    game1.fd.Close()
    game2.fd.Close()
    cleanup()
}

func (bt *bospec)tp()string {
    return runtime.FuncForPC(reflect.ValueOf(bt.boot).Pointer()).Name()
}

func doggleship() {
    go onexit()
    fmt.Printf("doggleship--salvo between two programs V7\n"+
        "Written 2022 by Eric Olson\n\n")
//    bt1:=bospec{bbcboot,"./bbcsalvo"}
//    bt1:=bospec{forboot,"./gfsalvo"}
//    bt2:=bospec{forboot,"./gfsalvo"}
//    bt1:=bospec{bbcboot,"./cheater"}
    bt1:=bospec{bbcboot,"./player1"}
    bt2:=bospec{bbcboot,"./player2"}
    fmt.Printf("Player one is %s of type %s\n",bt1.pname,bt1.tp())
    fmt.Printf("Player two is %s of type %s\n",bt2.pname,bt2.tp())
    fmt.Printf("\n%7s %7s %7s\n","Trial","Winner","Turns")
    for trial=1;trial<=tnum;trial++ {
        dotrial(bt1,bt2)
    }
    myexit(0)
}

func main(){
    flag.BoolVar(&batch,"b",false,"Only produce text output.")
    flag.IntVar(&tnum,"t",67,"Run specified number of trials.")
    flag.Parse()
    if batch {
        doggleship()
    } else {
        go doggleship()
        dogui()
    }
}
