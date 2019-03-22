libHana.a:
	go build -buildmode=c-archive -o libHana.a

clean:
	rm -rf hana
	rm -rf *.a *.h main *.a

