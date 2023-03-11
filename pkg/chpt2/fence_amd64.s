#include "textflag.h"

// func mfence()
TEXT ·mfence(SB), NOSPLIT, $0-0
MFENCE
RET
