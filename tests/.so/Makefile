
libfunction.so:
	gcc so.c --shared -o libfunction.so
libfunction_lin64.lzma: libfunction.so
	xz -z -9 -c libfunction.so > libfunction_lin64.lzma
libfunction_lin64.gzip: libfunction.so
	gzip -9 -c libfunction.so > libfunction_lin64.gzip
lin64_all: libfunction_lin64.gzip libfunction_lin64.lzma

function.dll:
	x86_64-w64-mingw32-gcc so.c --shared -o function.dll
libfunction_win64.lzma: function.dll
	xz -z -9 -c function.dll > libfunction_win64.lzma
libfunction_win64.gzip: function.dll
	gzip -9 -c function.dll > libfunction_win64.gzip
win64_all: libfunction_win64.gzip libfunction_win64.lzma

all: lin64_all win64_all

