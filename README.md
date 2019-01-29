# typograf — Go library for the ArtLebedev Studio Typograf SOAP webservice
The library for [ArtLebedev's Typograf](https://www.artlebedev.ru/typograf/about/).

There are a webserver and console utility in `cmd` folder. 

In `examples` you can find a JavaScript client example which utilizes the webserver from `cmd`.

## Console utility usage:

```
$ echo "- Это "Типограф"?" | typograf
```

Output:
```
<p>&#151;&nbsp;Это Типограф?<br />
</p>
```
