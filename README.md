# jplookup
Standalone CLI Japanese to English based on JMdict. This is a re-write of https://github.com/amyeo/jisho-cli mainly because of dependency frustrations with python.

# How to run

```
$ ./download_jmdict.sh
$ make index
$ make build
```

The first command downloads the JMdict dictionary from the source and unpacks it.
The second command makes the index for the dictionary

After running those two commands, you can use the dictionary

# Running

After running the commands above, the built executable will be at bin/jplookup.
```
$ bin/jplookup かいさつ

Results:
===================

ID: 1200840
Kanji (漢字):
改札

Kana (かな):
かいさつ

Meaning (ENG / 英語):
examination of tickets
ticket gate
ticket barrier

---------------------

ID: 1728410
Kanji (漢字):
開札

Kana (かな):
かいさつ

Meaning (ENG / 英語):
opening or unsealing of bids

(...)
```

You can then transfer the built executable into any directory that you desire.
** Take note that besides the executable, jisho.db (the built index) must exist in the same directory as the executable **

After this, you can delete "JMdict" and "JMdict.gz" if desired.

# Wildcard and placeholder queries

There are only 2 special queries in this program. Wildcards via '*' and '＊' and placeholders via '?' and '？'.

``` 観?車 ``` will match with ``` 観覧車 ``` and ``` 平*度 ``` matches with ``` 平均速度 ```

So you use '?' for a fixed/known length of unknown characters and '*' for an unknown length of unknown characters.

# Interactive and non-interactive queries

```
$ bin/jplookup かいさつ
```
will lookup a word, but if you do not specify a word:
```
$ bin/jplookup
Lookup > 気動車
```
It will prompt you for the word.

# JMdict attribution
This publication has included material from the JMdict (EDICT, etc.) dictionary files in accordance with the licence provisions of the Electronic Dictionaries Research Group. See http://www.edrdg.org/ 
