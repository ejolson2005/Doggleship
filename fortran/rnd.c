#include <stdint.h>
#include <fcntl.h>
#include <unistd.h>

typedef struct {
    uint64_t x,w,s;
} rstate;
static rstate gs;

static uint32_t rint32(rstate *p){
    p->x*=p->x; p->w+=p->s;
    p->x+=p->w; p->x=(p->x>>32)|(p->x<<32);
    return (uint32_t)p->x;
}

static void rseed() {
    int fd=open("/dev/urandom",O_RDONLY);
    gs.s=0xb5ad4eceda1ce2a9;
    if(fd>=0){
        read(fd,&gs,sizeof(gs));
        close(fd);
    }
    gs.s|=1;
}

float rnd_(float *x){
    if(gs.s==0) rseed();
    for(;;){
        float r=rint32(&gs)/4294967296.0;
        if(r<1.0) return r;
    }
}
