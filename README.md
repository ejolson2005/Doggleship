# Doggleship

## Salvo Between Programs

Another year has passed and the dog developer is having another birthday.
My plan is to create an optimal AI based on Bayesian statistics for playing
the battleship game described as SALVO in David Ahl's classic 101 Basic Computer Games where
Larry Siegel wrote
 
>SALVO is played on a 10x10 grid or board using an x,y coordinate system.
>The player has 4 ships: battleship (5 squares), cruiser (3 squares), and
>two destroyers (2 squares each). The ships must be placed horizontally,
>vertically, or diagonally and must not overlap. The ships do not move during the game.
>
>As long as any square of a battleship still survives,
>the player is allowed three shots, for a cruiser 2 shots, 
>and for each destroyer 1 shot. Thus, at the beginning of the game the player 
>has 3+2+1+1=7 shots. The player enters all of his shots and the computer 
>tells what was hit. A shot is entered by its grid coordinates, x,y. The
>winner is the one who sinks all of the opponent's ships. 

http://www.bitsavers.org/pdf/dec/_Books/101_BASIC_Computer_Games_Mar75.pdf

Note a Milton Bradly game with the same theme was marketed using the famous line
"you sank my battleship" as heard in the advertisement

https://www.youtube.com/watch?v=GkwMDkfrZ1M

This is a simplified version of SALVO. An even simpler version of the game was recently
finished in the Black Sea,
though it's not clear who sank the battleship.

In 2013 and 2014 a Pi versus Pi battleship tournament was part of
the Blue Pi Thinking program
at the University of York. In particular, Gordon Hollingworth wrote

>BattlePi is a game of battleships played automatically by the 
>Raspberry Pi with a server through which two Raspberry Pis can communicate.
>The students are given the (well commented Python) software which initially just 
>chooses random positions to make shots at. The student then has to modify
>the python code to implement better AI in both firing position and ship placement. 

See

https://www.raspberrypi.com/news/let-the-battlepi-commence/

and also

https://www.raspberrypi.com/news/blue-pi-thinking-from-the-university-of-york/

Before computing the Bayesian posterior after a salvo of Neptune cruise missiles,
the first part of Fido's new birthday project will be to strategically place the ships in
a statistically unbiased way and then develop a framework that allows the super pet to
play against other computers. I'm personally looking forward to
making an improved version of the program I wrote long ago for the Apple II.

This repository contains

-  Source for the original Basic program written by Larry Siegel with
   the typo introduced in the Basic Computer Games Microcomputer Edition
   corrected and adapted for BBC Basic.
   
-  A Fortran program which plays the same game and chooses the target
   based on a probability distribution calculated from all possible positions where the
   remaining ships might be.
   
-  A driver written in Go which allows either of the above programs to
   automatically play against each other.
   
-  A dashboard written in Go using the Ebiten gaming library which displays
   the status of a computer versus computer simulation.
   
## Building the Program

This program has been developed for Linux.
You need Go, gfortran and the console version of BBC Basic as well
as the prerequisites for Ebiten.

https://ebiten.org/

Compiling should be as simple as typing
```
$ make
```
in the main directory.  This will compile the Fortran code, tokenize the
Basic program and further produce a program called doggleship which by
default directs the Fortran program to play against itself.  The Basic
program can be specified by changing a single line in doggleship.go and
recompiling.

Note that the executable expects the Basic and Fortran programs that will
play against each other to be in the current working directory.

More information about this project appears in the Raspberry Pi Forum at

https://forums.raspberrypi.com/viewtopic.php?p=1994161#p1994161

Good Luck!
