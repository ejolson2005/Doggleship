package main

import(
    _ "embed"
    _ "image/png"
    _ "image/jpeg"
    "os"
    "fmt"
    "image/color"
    "bytes"
    "time"
    "strconv"
    "golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/text"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

//go:embed resources/DejaVuSansMono-Bold.ttf
var DejaVu []byte
    
//go:embed resources/background.jpg
var background []byte

var (
    tlimg,bgimg,shimg *ebiten.Image
    fntext font.Face
    fnhead font.Face
    count int=0
    bdleft,bdright [10][10]int
    dflag bool
)

const (
    xoff=32
    yoff=92
    tsize=32
    dpi=72
)

func (g *Game)Update()error {
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    dflag=true
    {
        op:=&ebiten.DrawImageOptions{}
        op.GeoM.Translate(0.0,0.0)
        screen.DrawImage(bgimg,op)
        text.Draw(screen,"Doggleship V4",fnhead,xoff,42,
            color.RGBA{0xff,0xff,0xff,0xff})
        gm:=""
        if dover>0 { gm=" [Game Over]" }
        text.Draw(screen,fmt.Sprintf("Trial %d Turn %d%s",trial,turn,gm),
            fntext,xoff,68,color.RGBA{0xff,0xff,0xff,0xff})
        co1:=color.RGBA{0xff,0xff,0xff,0xff}
        co2:=color.RGBA{0xff,0xff,0xff,0xff}
        if dover==1 {
            co1=color.RGBA{0x00,0xff,0x00,0xff}
            co2=color.RGBA{0xff,0x00,0x00,0xff}
        } else if dover==2 {
            co1=color.RGBA{0xff,0x00,0x00,0xff}
            co2=color.RGBA{0x00,0xff,0x00,0xff}
        }
        text.Draw(screen,fmt.Sprintf("Wins %d Losses %d",win1,win2),
            fntext,xoff,458,co1)
        text.Draw(screen,fmt.Sprintf("Wins %d Losses %d",win2,win1),
            fntext,xoff+13*(tsize+2),458,co2)
    }
    for i:=0;i<10;i++ {
        for j:=0;j<10;j++ {
            op:=&ebiten.DrawImageOptions{}
            x:=xoff+i*(tsize+2)
            y:=yoff+j*(tsize+2)
            op.GeoM.Translate(float64(x),float64(y))
            if ishere(&game2,i,j){
                screen.DrawImage(shimg,op)
            } else {
                screen.DrawImage(tlimg,op)
            }
            if bdleft[i][j]==0 { continue }
            str:=strconv.Itoa(bdleft[i][j])
            bound,_:=font.BoundString(fntext,str)
            w:=(bound.Max.X-bound.Min.X).Ceil()
            h:=(bound.Max.Y-bound.Min.Y).Ceil()
            x=x+(tsize-w)/2
            y=y+(tsize-h)/2 + h
            text.Draw(screen,str,fntext,x,y,
                color.RGBA{0x00,0x00,0x00,0xff})
        }
    }
    for i:=0;i<10;i++ {
        for j:=0;j<10;j++ {
            op:=&ebiten.DrawImageOptions{}
            x:=xoff+(i+13)*(tsize+2)
            y:=yoff+j*(tsize+2)
            op.GeoM.Translate(float64(x),float64(y))
            if ishere(&game1,i,j){
                screen.DrawImage(shimg,op)
            } else {
                screen.DrawImage(tlimg,op)
            }
            if bdright[i][j]==0 { continue }
            str:=strconv.Itoa(bdright[i][j])
            bound,_:=font.BoundString(fntext,str)
            w:=(bound.Max.X-bound.Min.X).Ceil()
            h:=(bound.Max.Y-bound.Min.Y).Ceil()
            x=x+(tsize-w)/2
            y=y+(tsize-h)/2 + h
            text.Draw(screen,str,fntext,x,y,
                color.RGBA{0x00,0x00,0x00,0xff})
        }
    }
}

func (g *Game)Layout(ow,oh int)(sw,sh int) {
    return 852,480
}

func doinit(){
    ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)
    ebiten.SetMaxTPS(1)
    tlimg=ebiten.NewImage(tsize,tsize)
    tlimg.Fill(color.RGBA{0x3f,0x7f,0x7f,0x2f})
    shimg=ebiten.NewImage(tsize,tsize)
    shimg.Fill(color.RGBA{0x6f,0x6f,0x6f,0x2f})
    var err error
    bgimg,_,err=ebitenutil.NewImageFromReader(
        bytes.NewReader(background))
    if err!=nil {
        fmt.Printf("Count not load gopher.png image!\n%v\n",err)
        os.Exit(1)
    }
    tt,err:=opentype.Parse(DejaVu)
    if err!=nil {
        fmt.Printf("%v\n",err)
        os.Exit(1)
    }
    fntext,err=opentype.NewFace(tt,&opentype.FaceOptions{
        Size: 18, DPI: dpi, Hinting: font.HintingFull})
    fnhead,err=opentype.NewFace(tt,&opentype.FaceOptions{
        Size: 24, DPI: dpi, Hinting: font.HintingFull})
}

func dorefresh(s time.Duration){
    if batch { return }
    dflag=false
    t:=twait
    for {
        ebiten.ScheduleFrame()
        time.Sleep(t*time.Millisecond)
        if dflag { break }
        if t<128 { t*=2 }
    }
    time.Sleep(s)
}

func dogui() {
    doinit()
    ebiten.SetWindowSize(852,480)
    ebiten.SetWindowTitle("doggleship")
    if err:=ebiten.RunGame(&Game{});err!=nil {
        fmt.Printf("Could not run doggleship dashboard!\n")
    }
}
