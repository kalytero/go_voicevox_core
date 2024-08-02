package main

/*
#cgo CFLAGS: -I${SRCDIR}/voicevox_core
#cgo LDFLAGS: -L${SRCDIR}/voicevox_core -l:libvoicevox_core.so -l:libonnxruntime.so.1.13.1

#include "voicevox_core.h"
#include <stdbool.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdio.h>
#include <fcntl.h>

bool CONST_TRUE = true;
*/
import "C"
import (
	"bufio"
	"fmt"
	"os"
	"time"
	"unsafe"
)

func main() {

	initialize_options := C.voicevox_make_default_initialize_options()
	initialize_options.load_all_models = C.CONST_TRUE
	initialize_options.open_jtalk_dict_dir = C.CString("voicevox_core/open_jtalk_dic_utf_8-1.11")
	fmt.Printf("[ OK ] Initialize Option\n")

	if C.voicevox_initialize(initialize_options) != C.VOICEVOX_RESULT_OK {
		fmt.Printf("[ FAILED ] Initialize Core\n")
		return
	}

	fmt.Printf("[ OK ] Initialize Core\n")

	for {
		var speaker_id C.uint32_t = 0
		var output_wav_length C.uintptr_t = 0
		var output_wav *C.uint8_t = nil
		var ch_res = make(chan C.int)
		var ch_wait = make(chan int)

		var sc = bufio.NewScanner(os.Stdin)
		fmt.Printf("> ")
		sc.Scan()
		var ss = sc.Text()
		
		go func() {
			var cs = C.CString(ss)
			ch_res <- C.voicevox_tts(cs, speaker_id, C.voicevox_make_default_tts_options(), &output_wav_length, &output_wav)
			C.free(unsafe.Pointer(cs))
			close(ch_res)
		}()

		go func() {
			for i := 0; i < 4; i++ {
				select {
				case <-ch_res:
					fmt.Printf("\r[ OK ] TTS\n")
					close(ch_wait)
					return
				default:
					spinChars := []rune{'|', '/', '-', '\\'}
					fmt.Printf("\r%c", spinChars[i])
					time.Sleep(10 * time.Millisecond)
					if i == 3 {
						i = 0
					}
				}
			}
		}()
		res := <-ch_res
		<-ch_wait

		if res != C.VOICEVOX_RESULT_OK {
			fmt.Printf("[ FAILED ] TTS\n")
			return
		} else {
			if fo, err := os.Create("output.wav"); err != nil {
				fmt.Printf("[ FAILED ] TTS\n")
				return
			} else {
				fo.Write(C.GoBytes(unsafe.Pointer(output_wav), C.int(output_wav_length)))
				fo.Close()
			}
		}

		if output_wav != nil {
			C.free(unsafe.Pointer(output_wav))
			output_wav = nil
		}
	}
}
