DEPS:=solana-ffi.h libsolana-ffi.a

all: $(DEPS)
.PHONY: all

$(DEPS): .install-solana-ffi  ;

.install-solana-ffi: rust
	cd rust && cargo build --release --all; cd ..
	find ./rust/target/release -type f -name "solana-ffi.h" -print0 | xargs -0 ls -t | head -n 1 | xargs -I {} cp {} ./cgo/solana-ffi.h
	find ./rust/target/release -type f -name "libsolana_ffi.a" -print0 | xargs -0 ls -t | head -n 1 | xargs -I {} cp {} ./cgo/libsolana_ffi.a
	c-for-go --ccincl solana-ffi.yml
	@touch $@

clean:
	rm -rf $(DEPS) .install-solana-ffi
	rm -rf cgo/*.go
	rm -rf cgo/*.h
	rm -rf cgo/*.a
.PHONY: clean
