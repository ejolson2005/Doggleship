/*  weyl32--Weyl sequence middle-square random number generator
    Written 2022 by Eric Olson */

#include <stdio.h>
#include "weyl32.h"

#ifdef HASU64

typedef long long unsigned int hwu64t;
typedef struct {
    hwu64t x,w,s;
} rstate;
static rstate gs;
static mlint rint32(rstate *p){
    p->x*=p->x; p->w+=p->s;
    p->x+=p->w; p->x=(p->x>>32)|(p->x<<32);
    return (mlint)p->x&0x7fffffff;
}

# ifdef DEBUG
static void outhex(hwu64t x){
    printf("%llX",x);
}
static void out64(hwu64t x){
    printf("%llu",x);
}
static mlint rint32d(rstate *p){
    p->x*=p->x; p->w+=p->s;
    printf("1: p->x=");out64(p->x);printf("\n");
    printf("2: p->w=");out64(p->w);printf("\n");
    p->x+=p->w;
    printf("3: p->x=");outhex(p->x);printf("\n");
    p->x=(p->x>>32)|(p->x<<32);
    printf("4: p->x=");outhex(p->x);printf("\n");
    return (mlint)p->x&0x7fffffff;
}
# endif

void rseed(mlint x){
    gs.x=x;
    gs.w=0;
    gs.s=13091206342165455529llu;
}

#else

typedef unsigned short muint;

#define my64len (8/(int)sizeof(muint))
#define my64hlf (4/(int)sizeof(muint))
#define my64msk ((muint)(-1))
#define my64sft (8*(int)sizeof(muint))

typedef muint my64[my64len];
typedef char mybcd[20];

static void add64(a,b) my64 a,b; {
    int i;
    mlint r,s=0;
    for(i=0;i<my64len;i++){
        r=(mlint)a[i]+b[i]+s;
        a[i]=r&my64msk;
        s=r>>my64sft;
    }
}

static void cadd64(a,c) my64 a; muint c; {
    int i;
    mlint r,s=c;
    for(i=0;i<my64len;i++){
        r=a[i]+s;
        a[i]=r&my64msk;
        s=r>>my64sft;
    }
}

static void cmul64(a,b) my64 a; muint b; {
    int i;
    mlint r,s=0;
    for(i=0;i<my64len;i++){
        r=(mlint)a[i]*b+s;
        a[i]=r&my64msk;
        s=r>>my64sft;
    }
}

static void mul2my64(a) my64 a; {
    int i;
    mlint r,s=0;
    for(i=0;i<my64len;i++){
        r=((mlint)a[i]<<1)+s;
        a[i]=r&my64msk;
        s=r>>my64sft;
    }
}

static void mul64(a,b) my64 a,b; {
    my64 c,d;
    int i;
    mlint r,s;
    for(i=0;i<my64len;i++) c[i]=a[i];
    for(i=0;i<my64len;i++) d[i]=0;
    for(i=0;i<my64len;i++){
        for(s=1;s<my64msk;s<<=1){
            r=b[i]&s;
            if(r) add64(d,c);
            mul2my64(c);
        }
    }
    for(i=0;i<my64len;i++) a[i]=d[i];
}

static int mydigit(a) char a; {
    if(a<'0'||a>'9') return 0;
    return 1;
}

# ifdef DEBUG

static void chexout(h) muint h; {
    if(h<10) putchar(h+'0');
    else putchar(h+'A'-10);
}

static void outbcd(a) mybcd a; {
    int i,f;
    f=0;
    i=19;
    for(;;){
        if(a[i]) f=1;
        if(f) putchar(a[i]+'0');
        if(i==0) break;
        i--;
    }
    if(f==0) putchar('0');
}

static void mul2bcd(a) mybcd a; {
    int i;
    unsigned int r,s=0;
    for(i=0;i<20;i++){
        r=(a[i]<<1)+s;
        if(r<10){
            a[i]=r;
            s=0;
        } else {
            a[i]=r-10;
            s=1;
        }
    }
}

static void incbcd(a) mybcd a; {
    int i,r,s=1;
    for(i=0;i<20;i++){
        r=a[i]+s;
        if(r<10){
            a[i]=r;
            s=0;
        } else {
            a[i]=r-10;
            s=1;
        }
    }
}

static void my64tobcd(c,a) mybcd c; my64 a; {
    int i;
    muint r,s;
    for(i=0;i<20;i++) c[i]=0;
    i=my64len-1;
    for(;;){
        s=1<<(my64sft-1);
        for(;;){
            r=a[i]&s;
            mul2bcd(c);
            if(r){
                incbcd(c);
            } 
            if(s==1) break;
            s>>=1;
        }
        if(i==0) break;
        i--;
    }
}

static void out64(a) my64 a; {
    mybcd c;
    my64tobcd(c,a);
    outbcd(c);
}

static void outhex(a) my64 a; {
    int i,j,f;
    muint r;
    f=0;
    i=my64len-1;
    for(;;){
        for(j=my64sft-4;j>=0;j-=4){
            r=(a[i]&((muint)0xf<<j))>>j;
            if(r) f=1;
            if(f) chexout(r);
        }
        if(i==0) break;
        i--;
    }
    if(f==0) chexout(0);
}

# endif

static void strtobcd(c,s) mybcd c; char *s; {
    int slen,i;
    for(slen=0;mydigit(s[slen]);slen++);
    for(i=0;i<slen;i++){
        c[i]=s[slen-i-1]-'0';
    }
    for(i=slen;i<20;i++) c[i]=0;
}

static void bcdtomy64(c,a) my64 c; mybcd a; {
    int i;
    for(i=0;i<my64len;i++) c[i]=0;
    i=19;
    for(;;){
        cmul64(c,10);
        cadd64(c,a[i]);
        if(i==0) break;
        i--;
    }
}

static void strtomy64(c,s) my64 c; char *s; {
    mybcd b;
    strtobcd(b,s);
    bcdtomy64(c,b);
}

static void mtswap(a) my64 a; {
    int i;
    muint r;
    for(i=0;i<my64hlf;i++){
        r=a[i]; a[i]=a[i+my64hlf]; a[i+my64hlf]=r;
    }
}

typedef struct {
    my64 x,w,s;
} rstate;

static rstate gs;
static mlint rint32(p) rstate *p; {
    int i;
    mlint r;
    mul64(p->x,p->x); add64(p->w,p->s);
    add64(p->x,p->w); mtswap(p->x);
    r=(my64msk>>1)&p->x[my64hlf-1];
    for(i=my64hlf-2;i>=0;i--) r=(r<<my64sft)+p->x[i];
    return r;
}

# ifdef DEBUG
static mlint rint32d(p) rstate *p; {
    int i;
    mlint r;
    mul64(p->x,p->x); add64(p->w,p->s);
    printf("1: p->x=");out64(p->x);printf("\n");
    printf("2: p->w=");out64(p->w);printf("\n");
    add64(p->x,p->w); 
    printf("3: p->x=");outhex(p->x);printf("\n");
    mtswap(p->x);
    printf("4: p->x=");outhex(p->x);printf("\n");
    r=(my64msk>>1)&p->x[my64hlf-1];
    for(i=my64hlf-2;i>=0;i--) r=(r<<my64sft)+p->x[i];
    return r;
}
# endif

void rseed(x) mlint x; {
    strtomy64(gs.x,my32toa(x));
    strtomy64(gs.w,"0");
    strtomy64(gs.s,"13091206342165455529");
}
#endif

char *my32toa(x) mlint x; {
    static char xb[128],*xs;
    if(xs-xb<22) xs=(&xb[128]);
    *--xs=0;
    if(x==0) *--xs='0';
    else while(x!=0){
        *--xs=x%10+'0';
        x/=10;
    }
    return xs;
}

unsigned int rdice(d) unsigned int d; {
    mlint r,p=0x7fffffff/d;
    for(;;){
        r=rint32(&gs)/p;
        if(r<d) return r;
    }
}

#ifdef DEBUG
mlint seeds[10]={0,1,2,3,7,11,13,17,19,23};
mlint sums[10]={
    1721296309, 741875526, 1761186872, 1326098620,
    926696101, 1562090973, 2100273949, 105764092,
    221130887, 1586030047 };
int main(){
    int i,j,k;
    mlint r;
    volatile mlint csum=0;
#ifndef HASU64
    if(sizeof(mlint)<=sizeof(muint)){
        printf("Error sizeof(mlint) not bigger than sizeof(luint)!\n");
        return 1;
    }
    printf("my64len=%d\n",my64len);
    printf("my64hlf=%d\n",my64hlf);
    printf("my64msk=%u\n",my64msk);
    printf("my64sft=%d\n",my64sft);
#endif
    rseed(12345);
    printf("gs.s="); outhex(gs.s);
    printf(" or "); out64(gs.s); printf("\n");
    for(i=0;i<4;i++){
        r=rint32d(&gs);
        printf("rint32(&gs)=%s\n",my32toa(r));
    }
    for(j=0;j<10;j++){
        rseed(seeds[j]);
        printf("gs.x="); outhex(gs.x);
        printf(" or "); out64(gs.x); printf("\n");
        for(i=0;i<10;i++){
            k=rdice(12);
            printf("%d ",k);
            csum=csum*13+k;
        }
        printf("\n");
        csum&=0x7fffffff;
        printf("Seed %s checksum %s %s validation\n",
            my32toa(seeds[j]),my32toa(csum),
            csum==sums[j]?"passed":"failed");
    }
    return 0;
}
#endif
