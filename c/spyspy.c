/*  spyspy.c--Spy Versus Spy Battleship V1
    Written 2022 by Eric Olson */

#include <stdio.h>
#include "weyl32.h"

#ifdef KANDR
extern char *malloc();
extern int exit();
#else
#include <unistd.h>
#include <stdlib.h>
#endif

#define N 10
static mlint gamma1=1900,gamma2=2200,rseq=127;
static int tmax=67;
static int trial,turn,quiet,winc,wind,verb;

typedef int bdspec[N][N];
static int dirx[8]={1,1,0,-1,-1,-1, 0, 1},
           diry[8]={0,1,1, 1, 0,-1,-1,-1};
static struct {
    char *name,*bd;
    int length,shots;
} ships[4] = {
    { "battleship"," B ",5,3 },
    { "cruiser"," C ",3,2 },
    { "destroyer<a>"," Da",2,1},
    { "destroyer<b>"," Db",2,1}};

static bdspec cat,dog;
static char *tframe,*bframe;
#ifdef ASCII
static char *petscii[8]={ "|","+","-","+","|","+","-","+" };
#else
#include "petscii.i"
#endif

#ifdef HASU64
static int shipfit(bdspec b,int x,int y,int d,int k);
static void shipput(bdspec b,int x,int y,int d,int k);
static void place(bdspec b);
static void pline(int z[]);
static char *dogadd(char *p,char *q);
static int tally(int f[4],bdspec b);
static int numshots(int f[4]);
static void getopen(int *x,int *y,bdspec b,int o);
static int fship(bdspec b,int s);
static char *getname(bdspec b);
static int attack(bdspec a,bdspec b,mlint gamma);
static char *fptoa(mlint x,int dp);
static mlint atofp(char *p,int dp);
extern int main(int argc,char *argv[]);
#endif

static int shipfit(b,x,y,d,k) bdspec b; int x,y,d,k; {
    int l,lnum=ships[k].length;
    for(l=0;l<lnum;l++){
        if(x<0||x>=N||y<0||y>=N) return 0;
        if(b[x][y]) return 0;
        x+=dirx[d]; y+=diry[d];
    }
    return 1;
}

static void shipput(b,x,y,d,k) bdspec b; int x,y,d,k; {
    int l,lnum=ships[k].length;
    k=(-k-1);
    for(l=0;l<lnum;l++){
        b[x][y]=k;
        x+=dirx[d]; y+=diry[d];
    }
}

static void place(b) bdspec b; {
    int i,j,k;
    for(i=0;i<N;i++) for(j=0;j<N;j++) b[i][j]=0;
    for(k=0;k<4;k++){
        for(;;){
            int x=rdice(N),y=rdice(N),d=rdice(8);
            if(shipfit(b,x,y,d,k)){
                shipput(b,x,y,d,k);
                break;
            }
        }
    }
}

static void pline(z) int z[]; {
    int j,k;
    for(j=0;j<N;j++){
        if(z[j]==0) {
            fputs("   ",stdout);
        } else if(z[j]<0) {
            k=(-z[j]-1);
            fputs(ships[k].bd,stdout);
        } else if(z[j]<10) {
            printf(" %1d ",z[j]);
        } else {
            printf(" %2d",z[j]);
        }
    }
}

static char *dogadd(p,q) char *p,*q; {
    while(*q) *p++=(*q++);
    return p;
}

static int doglen(p) char *p; {
    char *q=p;
    while(*q) q++;
    return q-p;
}

static void mkframe(){
    int i,k,tf,bf;
    char *p;
    if(tframe) return;
    tf=2*(doglen(petscii[3])+(3*N+1)*doglen(petscii[2])+doglen(petscii[1]));
    tframe=malloc(tf+4);
    p=tframe;
    for(k=0;k<2;k++){
        p=dogadd(p,petscii[3]);
        for(i=0;i<N;i++){
            p=dogadd(p,petscii[2]);
            p=dogadd(p,petscii[2]);
            p=dogadd(p,petscii[2]);
        }
        p=dogadd(p,petscii[2]);
        p=dogadd(p,petscii[1]);
        if(k==0) p=dogadd(p,"   ");
    }
    *p=0;
    bf=2*(doglen(petscii[5])+(3*N+1)*doglen(petscii[6])+doglen(petscii[7]));
    bframe=malloc(bf+4);
    p=bframe;
    for(k=0;k<2;k++){
        p=dogadd(p,petscii[5]);
        for(i=0;i<N;i++){
            p=dogadd(p,petscii[6]);
            p=dogadd(p,petscii[6]);
            p=dogadd(p,petscii[6]);
        }
        p=dogadd(p,petscii[6]);
        p=dogadd(p,petscii[7]);
        if(k==0) p=dogadd(p,"   ");
    }
    *p=0;
}

static void pboards(){
    int i;
    char buf[64];
    if(quiet) return;
    printf("Doggleship Spy Version 1\n");
    printf("Trial %d Turn %d\n",trial,turn);
    mkframe();
    puts(tframe);
    for(i=0;i<N;i++){
        printf("%s",petscii[4]);
        pline(dog[i]);
        printf(" %s   %s",petscii[0],petscii[4]);
        pline(cat[i]);
        printf(" %s\n",petscii[0]);
    }
    puts(bframe);
    sprintf(buf,"Wins %d Losses %d",winc,wind);
    printf("%-35s Wins %d Losses %d\n\n",buf,wind,winc);
}

static int tally(f,b) int f[4]; bdspec b; {
    int i,j,k,r;
    for(k=0;k<4;k++) f[k]=0;
    r=0;
    for(i=0;i<N;i++){
        for(j=0;j<N;j++){
            k=b[i][j];
            if(k<=0){
                r++;
                if(k<0){
                    f[-k-1]++;
                }
            }
        }
    }
    return r;
}

static int numshots(f) int f[4]; {
    int k,r;
    r=0;
    for(k=0;k<4;k++){
        if(f[k]) r+=ships[k].shots;
    }
    return r;
}

static void getopen(x,y,b,o) int *x,*y; bdspec b; int o; {
    int i,j;
    for(i=0;i<N;i++){
        for(j=0;j<N;j++){
            if(b[i][j]<=0){
                if(o==0){
                    *x=i; *y=j;
                    return;
                }
                o--;
            }
        }
    }
    fprintf(stderr,"Consistency error in getopen!\n");
    exit(1);
}

static int fship(b,s) bdspec b; int s; {
    int i,j,r,o;
    o=0;
    for(i=0;i<N;i++){
        for(j=0;j<N;j++){
            r=b[i][j];
            if(r<=0){
                if(r<0){
                    if(s==0) return o;
                    s--;
                }
                o++;
            }
        }
    }
    fprintf(stderr,"Consistency error in fship!\n");
    exit(1);
}

static char *getname(b) bdspec b; {
    if(b==cat){
        return "cat";
    } 
    return "dog";
}

static int attack(a,b,gamma) bdspec a,b; mlint gamma; {
    int f[4],s,snum,x[7],y[7],o[7],onum,oship,ocount;
    int i,j,k,fgood,g;
    char *aname,*bname;
    aname=getname(a); bname=getname(b);
    onum=tally(f,a);
    snum=numshots(f);
    onum=tally(f,b);
    if(snum==0){
        if(verb){
            printf("The %s has no more shots.\n",aname);
            printf("The %s has won!\n\n",bname);
        }
        return -1;
    }
    if(verb) printf("The %s has %d shots.\n",aname,snum);
    if(snum>onum){
        if(verb){
            printf("There are only %d remaining openings!\n",onum);
            printf("The %s has won!\n\n",aname);
        }
        return 1;
    }
    oship=0; ocount=0;
    for(k=0;k<4;k++) {
        oship+=f[k]; f[k]=0;
    }
    for(s=0;s<snum;s++){
        for(;;){
            fgood=1;
            g=rdice(10000);
            if(g<gamma&&ocount<oship){
                o[s]=fship(b,rdice(oship));
            } else {
                o[s]=rdice(onum);
            }
            for(i=0;i<s;i++){
                if(o[i]==o[s]) fgood=0;
            }
            if(fgood) break;
        }
        getopen(&i,&j,b,o[s]);
        if(verb) printf("Shot %d: (%d,%d)\n",s+1,i+1,j+1);
        k=(-b[i][j]-1);
        if(k>=0) { f[k]++; ocount++; }
        x[s]=i; y[s]=j;
    }
    for(s=0;s<snum;s++) b[x[s]][y[s]]=turn;
    for(k=0;k<4;k++){
        for(s=0;s<f[k];s++){
            if(verb){
                printf("The %s %s has been hit!\n",
                    bname,ships[k].name);
            }
        }
    }
    if(verb) putchar('\n');
    return 0;
}

static char *fptoa(x,dp) mlint x; int dp; {
    int d;
    static char xb[128],*xs;
    if(xs-xb<24) xs=(&xb[128]);
    *--xs=0;
    if(dp>0){
        for(d=0;d<dp;d++){
            *--xs=x%10+'0';
            x/=10;
        }
        *--xs='.';
    }
    if(x==0) *--xs='0';
    else while(x!=0){
        *--xs=x%10+'0';
        x/=10;
    }
    return xs;
}

static mlint atofp(p,dp) char *p; int dp; {
    mlint r=0;
    while(*p>='0'&&*p<='9'){
        r=r*10+*p-'0'; p++;
    }
    if(*p=='.') p++;
    while(dp>0){
        r*=10;
        if(*p>='0'&&*p<='9'){
            r+=(*p-'0'); p++;
        }
        dp--;
    }
    return r;
}

static void help(){
    fprintf(stderr,"Usage: spyspy [options]\n\n");
    fprintf(stderr,"where options are\n");
    fprintf(stderr,"\t-b  \tRun in batch mode without display.\n");
    fprintf(stderr,
        "\t-c x\tSet gamma to x percent for the cat (default %s).\n",
        fptoa(gamma1,2));
    fprintf(stderr,
        "\t-d x\tSet gamma to x percent for the dog (default %s).\n",
        fptoa(gamma2,2));
    fprintf(stderr,"\t-v  \tTurn on verbose mode.\n");
    fprintf(stderr,"\t-q  \tOnly print summary at the end.\n");
    fprintf(stderr,
        "\t-r n\tSet random seed to n (default %s).\n",fptoa(rseq,0));
    fprintf(stderr,
        "\t-t n\tRun n number of trials (default %d).\n",tmax);
    fprintf(stderr,"\t-h  \tPrint this message.\n");
    exit(1);
}

static void docmdline(argc,argv) int argc; char **argv; {
    while(--argc>0&&*argv[1]=='-'){
        argv++;
        while(*++*argv) switch(**argv){
case 'b':
            quiet=1;
            continue;
case 'c':
            if(*++*argv) gamma1=atofp(*argv,2);
            else if(--argc>0) gamma1=atofp(*++argv,2);
            else help();
            goto nextarg;
case 'd':
            if(*++*argv) gamma2=atofp(*argv,2);
            else if(--argc>0) gamma2=atofp(*++argv,2);
            else help();
            goto nextarg;
case 'q':
            quiet=2;
            continue;
case 'r':
            if(*++*argv) rseq=atofp(*argv,0);
            else if(--argc>0) rseq=atofp(*++argv,0);
            else help();
            goto nextarg;
case 't':
            if(*++*argv) tmax=atofp(*argv,0);
            else if(--argc>0) tmax=atofp(*++argv,0);
            else help();
            goto nextarg;
case 'v':
            verb=1;
            continue;
default:
            help();
        }
nextarg:;
    }
}

int main(argc,argv) int argc; char *argv[]; {
    int w;
    mlint pwin;
    printf("spyspy--Spy Versus Spy Battleship V1\n");
    printf("Written 2022 by Eric Olson\n\n");
    docmdline(argc,argv);
    rseed(rseq);
    printf("gamma1=%s, gamma2=%s, rseq=%s\n\n",
        fptoa(gamma1,2),fptoa(gamma2,2),fptoa(rseq,0));
    if(quiet==1){
        printf("%7s %7s %7s\n","Trial","Winner","Turns");
    }
    winc=0; wind=0;
    for(trial=1;trial<=tmax;trial++){
        place(cat);
        place(dog);
        turn=0;
        pboards();
        for(turn=1;;turn++){
            w=attack(cat,dog,gamma1);
            if(w) break;
            w=(-attack(dog,cat,gamma2));
            if(w) break;
            pboards();
        }
        if(w>0) winc++;
        else if(w<0) wind++;
        pboards();
        if(quiet==1) printf("%7d %7d %7d\n",trial,w>0?1:2,turn);
    }
    pwin=((mlint)winc*10000+(tmax>>1))/tmax;
    printf("Out of %d trials\n",tmax);
    printf("\t%s wins %s percent of the time.\n",
        getname(cat),fptoa(pwin,2));
    printf("\t%s wins %s percent of the time.\n",
        getname(dog),fptoa(10000-pwin,2));
    return 0;
}
