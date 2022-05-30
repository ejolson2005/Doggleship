#ifndef weyl32_h
#define weyl32_h

# ifdef KANDR
# define void int
# define volatile
# endif

# ifdef HASU64

typedef int mlint;
extern char *my32toa(mlint x);
extern void rseed(mlint x);
extern unsigned int rdice(unsigned int d);

# else

typedef long mlint;
extern char *my32toa();
extern void rseed();
extern unsigned int rdice();

# endif

#endif
