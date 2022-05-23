/*  doggleship.go--salvo between programs
    Written 2022 by Eric Olson

    Version 1 run SALVO from 101 Basic Computer Games
    Version 2 cleanup running programs on interrupt
    Version 3 add Fortran driver and refactor 
    Version 4 add Ebiten graphical dashboard
    Version 5 add option to turn of dashboard */

package main

import (
    "C"; "bufio"; "flag"; "fmt"; "os"; "os/exec"; "os/signal"
    "reflect"; "runtime"; "strconv"; "strings"; "syscall"
    "time"; "unsafe"
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
    gmspec struct {
        b []btspec
        sr []string
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
    return fmaster,cmd.Process
}

func monitor(sr *[]string,fp *os.File,log string) {
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
                (*sr)[len(*sr)-1]+=string(b[i:j])
                fmt.Fprintf(logio,"%s\n",b[i:j])
                *sr=append(*sr,"")
                i=j+1
            }
        }
        if i<j {
            (*sr)[len(*sr)-1]+=string(b[i:j])
            fmt.Fprintf(logio,"%s",b[i:j])
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
    if game.sr==nil {
        game.sr=make([]string,1,512)
    } else {
        game.sr=game.sr[:1]
        game.sr[0]=""
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

func (game *gmspec)findtail(s string)bool {
    t:=twait
    for p:=len(game.sr)-1;!istail(game.sr[p],s);p=len(game.sr)-1 {
        time.Sleep(t*time.Millisecond)
        if game.sr[p]==">" { return false }
        if p>2 {
            for q:=p-2;q<=p;q++ {
                if game.sr[q]=="I HAVE WON" { return false }
                if game.sr[q]=="YOU HAVE WON" { return false }
                if game.sr[q]==" YOU WIN!" { return false }
                if game.sr[q]==" YOU LOSE!" { return false }
            }
        }
        if t<128 { t*=2 }
    }
    return true
}

func (game *gmspec)nextprompt(p int,s string)int {
    t:=twait
    for {
        game.findtail(s)
        r:=len(game.sr)
        if r>p { return r }
        time.Sleep(t*time.Millisecond)
        if t<128 { t*=2 }
    }
}

func (game *gmspec)youhave()int{
    match:=func()int {
        p:=len(game.sr)
        for i:=p-3;i<p;i++ {
            if game.sr[i]=="I HAVE WON" { return i }
            if game.sr[i]=="YOU HAVE WON" { return i }
            if game.sr[i]==" YOU WIN!" { return i }
            if game.sr[i]==" YOU LOSE!" { return i }
            if strings.Contains(game.sr[i],"AND YOU HAVE") { return i }
        }
        return 0
    }
    t:=twait
    for {
        p:=match()
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
    game.findtail("DO YOU WANT TO START? ")
    p:=len(game.sr)
    fmt.Fprintf(game.fd,"WHERE ARE YOUR SHIPS?\n")
    game.nextprompt(p,"? ")
    j:=-1
    for i:=p;i<len(game.sr);i++ {
        r:=getname(game.sr[i])
        if r>=0 { j=r
        } else if j>=0{
            sxy:=strings.Fields(game.sr[i])
            if len(sxy)==2 {
                x,err:=strconv.Atoi(sxy[0])
                if err!=nil { continue }
                y,err:=strconv.Atoi(sxy[1])
                if err!=nil { continue }
                game.b[j].p=append(game.b[j].p,cospec{x,y})
            }
        }
    }
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
    game.findtail("CHOICE: ")
    p:=len(game.sr)
    fmt.Fprintf(game.fd,"0\n")
    game.nextprompt(p,": ")
    var (them [10][10]int; q=0)
    for i:=p;i<len(game.sr);i++ {
        if q>0{
            if istail(game.sr[i],"-----"){ break }
            xs:=strings.Fields(game.sr[i])
            for j:=0;j<len(xs);j++ {
                them[i-q][j],_=strconv.Atoi(xs[j])
            }
        }
        if istail(game.sr[i],"-----"){ 
            q=i+1
        }
    }
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
    p:=len(game.sr)-1
    for i:=0;i<len(b);i++ {
        for j:=0;j<len(b[i].p);j++ {
            p=game.nextprompt(p,"? ")
            fmt.Fprintf(game.fd,"%d,%d\n",
                b[i].p[j].i,b[i].p[j].j)
        }
    }
    game.findtail("DO YOU WANT TO SEE MY SHOTS? ")
    p=len(game.sr)
    fmt.Fprintf(game.fd,"YES\n")
    game.nextprompt(p,"? ")
}

var dirs=[8]cospec{
    cospec{0,1},cospec{-1,1},cospec{-1,0},cospec{-1,-1},
    cospec{0,-1},cospec{1,-1},cospec{1,0},cospec{1,1}}

func (game *gmspec)forplace(b []btspec){
    p:=len(game.sr)-1
    for i:=0;i<len(b);i++ {
        p=game.nextprompt(p,": ")
        fmt.Fprintf(game.fd,"%d%d\n",
            b[i].p[0].i-1,b[i].p[0].j-1)
        p=game.nextprompt(p,": ")
        d:=cospec{b[i].p[1].i-b[i].p[0].i,b[i].p[1].j-b[i].p[0].j}
        var j int
        for j=0;j<len(dirs);j++ {
            if dirs[j]==d { break }
        }
           fmt.Fprintf(game.fd,"%d\n",j+1)
    }
    p=game.youhave()
    game.nextprompt(p,": ")
}

func (game *gmspec)bbcoutgoing(r []cospec){
    p:=len(game.sr)-1
    for i:=0;i<len(r);i++ {
        p=game.nextprompt(p,"? ")
        fmt.Fprintf(game.fd,"%d,%d\n",r[i].i,r[i].j)
    }
    game.nextprompt(p,"? ")
}

func (game *gmspec)foroutgoing(r []cospec){
    p:=len(game.sr)-1
    for i:=0;i<len(r);i++ {
        p=game.nextprompt(p,": ")
        fmt.Fprintf(game.fd,"%d%d\n",r[i].i-1,r[i].j-1)
    }
    game.youhave()
    game.nextprompt(p,": ")
}

func (game *gmspec)bbcincoming()([]cospec,stspec) {
    r:=make([]cospec,0,7)
    var p int
    for p=len(game.sr)-1;p>=0;p-- {
        if ishead(game.sr[p],"I HAVE WON"){
            return r,won
        } else if ishead(game.sr[p],"YOU HAVE WON"){
            return r,lost
        } else if ishead(game.sr[p],"I HAVE "){
            break
        }
    }
    sn:=strings.Fields(game.sr[p])
    if len(sn)!=4 {
        fmt.Printf("%s:Couldn't parse number of shots from '%s'\n",
            game.log,game.sr[p])
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
        sxy:=strings.Fields(game.sr[i])
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
    var p int
    for p=len(game.sr)-1;p>=0;p-- {
        if ishead(game.sr[p]," YOU LOSE!"){
            return r,won
        } else if ishead(game.sr[p]," YOU WIN!"){
            return r,lost
        } else if strings.Contains(game.sr[p],"AND I HAVE"){
            break
        }
    }
    sn:=strings.Fields(game.sr[p])
    if len(sn)!=9 {
        fmt.Printf("%s:Couldn't parse number of shots from '%s'\n",
            game.log,game.sr[p])
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
        sxy:=strings.Fields(game.sr[i])
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
    fmt.Printf("doggleship--salvo between two programs V5\n"+
        "Written 2022 by Eric Olson\n\n")
//    bt1:=bospec{bbcboot,"./bbcsalvo"}
//    bt1:=bospec{forboot,"./gfsalvo"}
//    bt2:=bospec{forboot,"./gfsalvo"}
    bt1:=bospec{bbcboot,"./cheater"}
    bt2:=bospec{bbcboot,"./cheater"}
    fmt.Printf("Player one is %s of type %s\n",bt1.pname,bt1.tp())
    fmt.Printf("Player two is %s of type %s\n",bt1.pname,bt2.tp())
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
