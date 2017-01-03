#!/bin/sh
HOST=localhost
PORT=12345

func todoadd()
{
    items="";
    for var in "$@"
    do
	items="$items -d item=\"$var\""
    done
    items="curl $HOST:$PORT/todo/add $items"
    eval $items
}
