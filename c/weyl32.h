#ifndef weyl32_h
#define weyl32_h

# ifdef HASU64

#include <stdint.h>
typedef int32_t mint32;
typedef uint8_t muint8;
extern char *my32toa(mint32 x);
extern void rseed(mint32 x);
extern unsigned int rdice(unsigned int d);

# else

typedef long mint32;
typedef unsigned int muint8;
extern char *my32toa();
extern void rseed();
extern unsigned int rdice();

# endif

#endif
