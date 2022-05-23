/*  cheater--salvo using covert military intelligence
    Written 2022 by Eric Olson

    Version 1 play SALVO with random cheating */

import Glibc

let π=4*atan(1.0)
let N=10
let γ=0.1

func getinput()->String {
    print("? ",terminator:"")
    let result=readLine()
    if result==nil {
        exit(1)
    } 
    return result!.uppercased()
}

struct coordsp {
    var x=0,y=0
}
struct boatsp {
    var s="",l=0,h=0
}
typealias gridsp=[[Int]]

var boats=[
    boatsp(s:"BATTLESHIP",l:5,h:3),
    boatsp(s:"CRUISER",l:3,h:2),
    boatsp(s:"DESTROYER<A>",l:2,h:1),
    boatsp(s:"DESTROYER<B>",l:2,h:1)]
var Δ=Array(repeating:coordsp(),count:8)
for d in 0...7 {
    Δ[d].x=Int(cos(Double(d)*π/4)*1.5)
    Δ[d].y=Int(sin(Double(d)*π/4)*1.5)
}
var noisy:Bool=true

func inrange(_ c:coordsp)->Bool{
    if c.x<0||c.x>=N||c.y<0||c.y>=N {
        return false
    }
    return true
}

func doesfit(_ 🐮:inout gridsp,
    _ xy:coordsp,_ d:Int,_ b:Int)->Bool{
    var c=xy
    for _ in 0..<boats[b].l {
        if 🐮[c.x][c.y]<0 {
            return false
        }
        c.x+=Δ[d].x; c.y+=Δ[d].y
        if !inrange(c) {
            return false
        }
    }
    return true
}

func placeit(_ 🐮:inout gridsp,
    _ xy:coordsp,_ d:Int,_ b:Int){
    var c=xy
    for _ in 0..<boats[b].l {
        🐮[c.x][c.y]=(-b-1)
        c.x+=Δ[d].x; c.y+=Δ[d].y
    }
}

func prgrid(_ 🐮:inout gridsp){
    for x in 0..<N {
        for y in 0..<N {
            print(🐮[x][y],terminator:" ")
        }
        print()
    }
}

func putships(_ 🐮:inout gridsp){
    for b in 0..<boats.count {
        while true {
            let xy=coordsp(x:Int.random(in:0..<N),
                y:Int.random(in:0..<N))
            let d=Int.random(in:0...7)
            if doesfit(&🐮,xy,d,b) {
                placeit(&🐮,xy,d,b)
                break
            }
        }
    }

}

func prships(_ 🐮:inout gridsp){
    for b in 0..<boats.count {
        print(boats[b].s)
        for x in 0..<N {
            for y in 0..<N {
                if 🐮[x][y]==(-b-1) {
                    print(" \(x+1)  \(y+1)")
                }
            }
        }
    }
}

func prlegal(){
    print("ILLEGAL, ENTER AGAIN")
}

func xyinput()->coordsp{
    while true {
        let s=getinput()
        let xy=s.split(separator:",")
        if xy.count==2 {    
            let x=Int(xy[0]),y=Int(xy[1])
            if !(x==nil||y==nil){
                return coordsp(x:x!-1,y:y!-1)
            } else {
                prlegal()
            }
        } else {
            prlegal()
        }
    }
}

func getships(_ 🐮:inout gridsp){
    for b in 0..<boats.count {
        print(boats[b].s)
        for _ in 0..<boats[b].l {
            let c={()->coordsp in
                while true {
                    let c=xyinput()
                    if inrange(c) {
                        let β=(-🐮[c.x][c.y]-1)
                        if β<0 {
                            return c
                        } else {
                            print("A",boats[β].s,"IS ALREADY THERE")
                        }
                    } else {
                        prlegal()
                    }
                }
            }()
            🐮[c.x][c.y]=(-b-1)
        }
    }
}

func numshots(_ 🐮:inout gridsp)->Int {
    var sunk=Array(repeating:true,count:boats.count)
    for x in 0..<N {
        for y in 0..<N {
            let β=(-🐮[x][y]-1)
            if β>=0 {
                sunk[β]=false
            }
        }
    }
    var r=0
    for b in 0..<boats.count {
        if !sunk[b] {
            r+=boats[b].h
        }
    }
    return r
}

func numopen(_ 🐮:inout gridsp)->Int {
    var r=0
    for x in 0..<N {
        for y in 0..<N {
            if 🐮[x][y]<=0 {
                r+=1
            }
        }
    }
    return r
}

func dowin(){
    print("I HAVE WON")
    exit(0)
}

func dolose(){
    print("YOU HAVE WON")
    exit(0)
}

func getempty(_ 🐮:inout gridsp)->coordsp{
    let ρ={()->Double in
        if numshots(&🐮)==0 {
            return 1.0
        }
        return Double.random(in:0..<1.0)
    }()
    while true {
        let xy=coordsp(x:Int.random(in:0..<N),
            y:Int.random(in:0..<N))
        let μ=🐮[xy.x][xy.y]
        if μ<=0 {
            if ρ>=γ||μ<0 {
                return xy
            }
        }
    }
}

func outgoing(_ 🐶:inout gridsp,_ 🐱:inout gridsp,_ t:Int){
    let m=numopen(&🐱)
    let n=numshots(&🐶)
    print("I HAVE \(n) SHOTS")
    if m<n {
        print("THE NUMBER OF MY SHOTS EXCEEDS",
            "THE NUMBER OF BLANK SQUARES")
        dowin()
    }
    var hits=Array(repeating:0,count:boats.count)
    for _ in 0..<n {
        let c=getempty(&🐱)
        let β=(-🐱[c.x][c.y]-1)
        if β>=0 {
            hits[β]+=1
        }
        🐱[c.x][c.y]=t
        print(" \(c.x+1)  \(c.y+1)")
    }
    for b in 0..<boats.count {
        for _ in 0..<hits[b] {
            print("I HIT YOUR",boats[b].s)
        }
    }
    let o=numshots(&🐱)
    if o==0 {
        dowin()
    }
}

func incoming(_ 🐱:inout gridsp,_ 🐶:inout gridsp,_ t:Int){
    let m=numopen(&🐶)
    let n=numshots(&🐱)
    print("YOU HAVE \(n) SHOTS")
    if m<n {
        print("THE NUMBER OF YOUR SHOTS EXCEEDS",
            "THE NUMBER OF BLANK SQUARES")
        dolose()
    }
    var hits=Array(repeating:0,count:boats.count)
    for _ in 0..<n {
        let c={()->coordsp in
            while true {
                let c=xyinput()
                if inrange(c) {
                    let τ=🐶[c.x][c.y]
                    if τ>0 {
                        print("YOU SHOT THERE BEFORE ON TURN",τ)
                    } else {
                        return c
                    }
                } else {
                    prlegal()
                }
            }
        }()
        let β=(-🐶[c.x][c.y]-1)
        if β>=0 {
            hits[β]+=1
        }
        🐶[c.x][c.y]=t
    }
    for b in 0..<boats.count {
        for _ in 0..<hits[b] {
            print("YOU HIT MY",boats[b].s)
        }
    }
    let o=numshots(&🐶)
    if o==0 {
        dolose()
    }    
}

func main(){
    print("cheater--salvo using covert military intelligence V1",
        "\nWritten 2022 by Eric Olson\n")
    var 🐶=Array(repeating:Array(repeating:0,count:N),count:N) as gridsp
    putships(&🐶)
    let start={ ()->String in
        while true {
            print("DO YOU WANT TO START",terminator:"")
            let start=getinput()
            if start.hasPrefix("W") {
                prships(&🐶)
            } else {
                return start
            }
        }
    }()
    print("ENTER COORDINATES FOR...")
    var 🐱=Array(repeating:Array(repeating:0,count:N),count:N) as gridsp
    getships(&🐱)
    print("DO YOU WANT TO SEE MY SHOTS",terminator:"")
    noisy=getinput().hasPrefix("Y")
    var t=1
    while true {
        print("\nTURN \(t)")
        if start.hasPrefix("Y"){
            incoming(&🐱,&🐶,t)
            outgoing(&🐶,&🐱,t)
        } else {
            outgoing(&🐶,&🐱,t)
            incoming(&🐱,&🐶,t)
        }
        t+=1
    }
}

main()
