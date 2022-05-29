/*  Indexing of petscii is for following positions of box frame.

        3 2 1
        4   0
        5 6 7

    Use strings to allow inclusion of utf8 and escape characters.  */

static char *petscii[8]={ "┃","┓","━","┏","┃","┗","━","┛" };
