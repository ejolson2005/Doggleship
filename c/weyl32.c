/*  weyl32--Weyl sequence middle-square random number generator
    Written 2022 by Eric Olson */

#include <stdio.h>
#include <stdint.h>

//#define DEBUG

//  Note cc65 doesn't support unsigned 64-bit integer types
//#define HASU64

char *u32toa(uint32_t x){
	static char xb[128],*xs;
	if(xs-xb<22) xs=&xb[128];
	*--xs=0;
	if(x==0) *--xs='0';
	else while(x>0){
		*--xs=x%10+'0';
		x/=10;
	}
	return xs;
}

#ifdef HASU64

typedef struct {
    uint64_t x,w,s;
} rstate;
static rstate gs;
static uint32_t rint32(rstate *p){
    p->x*=p->x; p->w+=p->s;
    p->x+=p->w; p->x=(p->x>>32)|(p->x<<32);
    return (uint32_t)p->x;
}
void rseed(unsigned int x){
	gs.x=x;
	gs.w=0;
    gs.s=13091206342165455529u;
}

#else

typedef unsigned char my64[8];
typedef unsigned char mybcd[20];

static void add64(my64 a, my64 b){
	int i;
	unsigned int r,s=0;
	for(i=0;i<8;i++){
		r=a[i]+b[i]+s;
		a[i]=r&0xff;
		s=r>>8;
	}
}

static void cadd64(my64 a, unsigned c){
	int i;
	unsigned int r,s=c;
	for(i=0;i<8;i++){
		r=a[i]+s;
		a[i]=r&0xff;
		s=r>>8;
	}
}

static void cmul64(my64 a,unsigned char b){
	int i;
	unsigned int r,s=0;
	for(i=0;i<8;i++){
		r=a[i]*b+s;
		a[i]=r&0xff;
		s=r>>8;
	}
}

static void mul2my64(my64 a){
	int i;
	unsigned int r,s=0;
	for(i=0;i<8;i++){
		r=(a[i]<<1)+s;
		a[i]=r&0xff;
		s=r>>8;
	}
}

static void mul64(my64 a,my64 b){
	my64 c,d;
	int i,j;
	unsigned int r;
	for(i=0;i<8;i++) c[i]=a[i];
	for(i=0;i<8;i++) d[i]=0;
	for(i=0;i<8;i++){
		for(j=1;j<256;j<<=1){
			r=b[i]&j;
			if(r) add64(d,c);
			mul2my64(c);
		}
	}
	for(i=0;i<8;i++) a[i]=d[i];
}

static int mydigit(char a){
	if(a<'0'||a>'9') return 0;
	return 1;
}

# ifdef DPRINT

static void outbcd(mybcd a){
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

static void my64tobcd(mybcd c,my64 a){
	int i,j;
	unsigned int r;
	for(i=0;i<20;i++) c[i]=0;
	i=7;
	for(;;){
		j=128;
		for(;;){
			r=a[i]&j;
			mul2bcd(c);
			if(r){
				incbcd(c);
			} 
			if(j==1) break;
			j>>=1;
		}
		if(i==0) break;
		i--;
	}
}

static void incbcd(mybcd a){
	int i;
	unsigned int r,s=1;
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

static void mul2bcd(mybcd a){
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

static void out64(my64 a){
	mybcd c;
	my64tobcd(c,a);
	outbcd(c);
}

static void chexout(unsigned char h){
	if(h<10) putchar(h+'0');
	else putchar(h+'a'-10);
}

static void outhex(my64 a){
	int i,f;
	unsigned int r;
	f=0;
	i=7;
	for(;;){
		r=(a[i]&0xf0)>>4;
		if(r) f=1;
		if(f) chexout(r);
		r=a[i]&0xf;
		if(r) f=1;
		if(f) chexout(r);
		if(i==0) break;
		i--;
	}
}

# endif

static void strtobcd(mybcd c,char *s){
	int slen,i;
	for(slen=0;mydigit(s[slen]);slen++);
	for(i=0;i<slen;i++){
		c[i]=s[slen-i-1]-'0';
	}
	for(i=slen;i<20;i++) c[i]=0;
}

static void bcdtomy64(my64 c,mybcd a){
	int i;
	for(i=0;i<8;i++) c[i]=0;
	i=19;
	for(;;){
		cmul64(c,10);
		cadd64(c,a[i]);
		if(i==0) break;
		i--;
	}
}

static void strtomy64(my64 c,char *s){
	mybcd b;
	strtobcd(b,s);
	bcdtomy64(c,b);
}

static void mtswap(my64 a){
	int i;
	unsigned int r;
	for(i=0;i<4;i++){
		r=a[i]; a[i]=a[i+4]; a[i+4]=r;
	}
}

typedef struct {
	my64 x,w,s;
} rstate;
static rstate gs;
static uint32_t rint32(rstate *p){
	mul64(p->x,p->x); add64(p->w,p->s);
	add64(p->x,p->w); mtswap(p->x);
	return *(uint32_t *)p->x;
}
void rseed(uint32_t x){
	strtomy64(gs.x,u32toa(x));
	strtomy64(gs.w,"0");
	strtomy64(gs.s,"13091206342165455529");
}
#endif

unsigned int rdice(unsigned int d){
	uint32_t p=4294967295/d;
	for(;;){
		uint32_t r=rint32(&gs)/p;
		if(r<d) return r;
	}
}

#ifdef DEBUG
int main(){
	int i,j;
	uint32_t seeds[10]={0,1,2,3,7,11,13,17,19,23};
	uint32_t sums[10]={
		2905269088u, 1437847190u, 3851704249u, 513246034u,
		3058933691u, 847856151u, 1738049419u, 3663490177u,
		1406267573u, 1714902782u };
	for(j=0;j<10;j++){
		rseed(seeds[j]);
		uint32_t csum=0;
		for(i=0;i<10;i++){
			int r=rdice(12);
			csum=csum*13+r;
		}
		printf("Seed %s checksum %s %s validation\n",
			u32toa(seeds[j]),
			u32toa(csum),csum==sums[j]?"passed":"failed");
	}
	return 0;
}
#endif
