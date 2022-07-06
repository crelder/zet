# FAQ

## General

1. Is there a tutorial?

> `zet init example` will download a small example zettelkasten and explains how to use it.

2. How do I add new zettel?

> You write your thought in a textfile or on a slip of paper (and scan it). Then you write the filename by hand and move/copy the file into the folder `zettelkasten/zettel`.

> For textfiles that have a predefined header format there is also the option of creating the filename automatically. `zet init example` will download a tutorial where the header format is explained.

3. How do I know if my zettelkasten is valid?

> `zet validate` will check for all kinds of inconsistencies like double ids or dead links in your zettelkasten.

4. What is the difference between "keywords" and "index topics"?

> The central concept is the index with entry points into your line of thoughts regarding a specific topic. Keywords are just a helper to already give some value via finding specific thoughts that you remember and therefor bridging the time between your first zettel and having enough zettel so that the power of the index will become visible.

5. What is the process of linking thoughts?

> First I write my thought down. Then I use the index for evaluating existing chains of thoughts and if I can link a thought there. If there is no possible predecessor for my thought to link to I use `VIEWS/unlinked/` to view suggestions of similar zettel.

6. How many zettel should be listed per thematic entry point (index entry)?

> Luhmann had maximum four IDs per index entry.
> "Auffällig ist, dass in der Regel bei einem Schlagwort maximal vier Verweisstellen im Kasten aufgeführt waren, also kein Anspruch auf vollständige Erfassung aller für das jeweilige Schlagwort relevanten Zettel erhoben wurde." [Source](https://niklas-luhmann-archiv.de/nachlass/zettelkasten).

7. Which tools do you combine `zet` with?

> This cli tool has its full potential when combined with other tools.

> * Syncing across all your devices: e.g., via iCloud, Nextcloud, rsync and a cronjob.
> * Version control: e.g., via Git
> * Scanning paper zettel: e.g., via the Adobe Scan App
> * Dictate thoughts: e.g., IA Writer with ios build-in speech recognition

> These tools will most likely change over your lifetime. The **invariant** of this system is the three zettelkasten elements: 1. text and images with your thought, 2. a text file with an index and 3. a text file with  literature references.

8. Is the zettelkasten a wiki?

> No. There are at least these two differences:
> 1. A wiki creates a network of entities, a zettelkasten a tree of entities.
> 2. The entities, the zettel, in the zettelkasten are immutable, in a wiki the entities, the pages, are mutable. A zettel contains a thought, a fact, which doesn't change. If your opinion has changed about this thought, create a follow-up zettel to this one.

9. What is the inspiration for `zet`?

> Niklas Luhmann's zettelkasten. I recommend the following content for learning about his zettelkasten:

> * [Niklas Luhmann's Zettelkasten digitalized](http://ds.ub.uni-bielefeld.de/viewer/collections/zettelkasten/)
> * A short essay from Luhmann on how to read: [Niklas Luhmann. Lesen Lernen](https://media.suhrkamp.de/mediadelivery/asset/cf9bb33d79fa476095a84f881aa0ca59/schriften-zu-kunst-und-literatur_9783518294727_leseprobe.pdf?contentdisposition=inline). Suhrkamp. 2008.
> * A short essay from Luhmann on how to work with his zettelkasten: [Kommunikation mit Zettelkästen](https://ckrybus.com/static/papers/luhmann1981.pdf)
> * [Niklas Luhmanns describes his zettelkasten](https://www.youtube.com/watch?v=mCFP5i_0ibE)
> * A [presentation](https://strengejacke.files.wordpress.com/2015/10/introduction-into-luhmanns-zettelkasten-thinking.pdf) on the questions how Niklas Luhmann worked with his zettelkasten 
> * A blog entry on [how to read](https://strengejacke.wordpress.com/2007/08/04/lesen-lernen/)
> * [Johannes Schmidt: Der Zettelkasten als Zweitgedächtnis Niklas Luhmanns, 30. April 2016, Kunstverein Hannover](https://vimeo.com/173128404)
> * [Einblicke in das System der Zettel - Geheimnis um Niklas Luhmanns Zettelkasten](https://www.youtube.com/watch?v=4veq2i3teVk)

10. What are alternatives to this project?

> There are alternatives for personal knowledge management systems like [infinitymaps](https://infinitymaps.io/imapping-tool/), wikis, [obsidian](https://obsidian.md/), [devonthink](https://www.devontechnologies.com/de/apps/devonthink), [orgmode](https://orgmode.org/)

> An alternative zettelkasten implementation is [zkn3](http://zettelkasten.danielluedecke.de/).

11. Is the zettelkasten just for collecting thoughts in an associative manner?

> No. Using the index and your chain of thoughts you can also use it to develop an own work about a topic (e.g. a paper, a presentation). If you link your thoughts with `zet views`, they will be put in a local order in folder `VIEWS/index`.

12. Is the zettelkasten only for raw thoughts?

> In the folder `zettel` are your raw thoughts. If you want, you can add a folder `works/` in your zettelkasten, which contains your finished works like book chapters, thesis, presentations, etc. You can name the filenames like you name your zettel. The program will then consider two folders when scanning: `zettel` and `works`.

> But actually I am stil not sure if this extra folder is necessary. you could as well include your own works in the zettelkasten just by referencing them like any other work via a bibkey and an entry in the `references.bib`.

13. What is the estimated maximum amount of thoughts one will enter over a lifetime?

> Around 100.000. Niklas Luhmann wrote in 45 years of working full-time as a knowledge worker (professor), around 90.000 handwritten paper slips ("zettel") with own thoughts.

14. Which size does the paper slips have?

> Luhmann used DIN A6 slipes of paper in his zettelkasten. This is also what I use, because it is a good trade of enough space for getting a concrete thought down and at the same time not too big, so you are protected from writing down several thoughts on one zettel. Also, I believe that limitation can trigger creativity and thinking.

15. How can I link an idea when I don't want to create a link to a predecessor?

> You can write "regarding this thought here, see also 190212f". But it is a lower priority link between two zettel than the Folgezettel / Predecessor.

16. Is the logical structure of the `VIEWS/index` created by `zet views` exactly the physical structure Luhmann had in his
    zettelkasten?

> No, since Luhmann could from one thought only branch of maximum two other thoughts. In zet you can branch of as many as you like.

17. When should I not consider using a zettelkasten approach?

> If you have to deliver a specific project with a deadline (e.g., thesis, dissertation, book, essay), then it is probably better to fill a file directly with your thoughts (instead of filling the zettelkasten first and then extracting relevant ideas out of it). Because only after a certain amount of thoughts in your system, the system gets really valuable so that you can extract "interesting" combinations of knowledge that suprise you.

## About this project


1. Is this project finished?

> This project is still a work in progress before reaching version 1. The core features are validation, generating entry points in your zettelkasten (views) and supporting the import of new zettel in text form. They are all available via `zet validate`, `zet views` and `zet import <PATH>`. See current tasks and planned features [here](./TODO.md).
> I might add a gephi file export and a browser UI for looking at the folgezettel structure, but I see this as optional.

