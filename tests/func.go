package tests

/*
int function(int);

#define DEFINE_JUMPER(x) \
        void *_godl_##x = (void*)0; \
        __asm__(".global "#x"\n\t"#x":\n\tmovq _godl_"#x"(%rip),%rax\n\tjmp *%rax\n")

DEFINE_JUMPER(function);
*/
import "C"
import "unsafe"

func function(i int) int {
	return int(C.function(C.int(i)))
}

func functionPtr() unsafe.Pointer {
	return unsafe.Pointer(&C._godl_function)
}
