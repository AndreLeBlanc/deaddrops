# Deadrop

**OSPP (1DT096) 2016 - Grupp 03**

Det här projektet ämnar att underlätta temporär filöverföring mellan två eller flera parter. En användare kan ladda upp en eller flera filer samt ange livslängd i form av dagar eller antal nedladdningar. Därefter så sparas de under en hashad URL och väntar på besökare. Servern som Deadrop körs på är skriven i Go och front-end i Javascript.

## Installera

### Back-end

##### Steg 1, installera Go
För att installera Go och se till så att ens GOROOT och GOPATH är rätt konfigurerade, rekommenderas det att följa
[Golangs officiella installationsguide](https://golang.org/doc/install).

##### Steg 2, klona projektet
Se till att stå i `$GOPATH/src` och skriv sen:

`git clone https://github.com/uu-it-teaching/ospp-2016-group-03/`

Därefter gå in i den nya katalogen och skriv:

`make dep` TODO: Fix this rule in the Makefile (make dependencies), change name of project folder to Deadrop and install database, etc.

### Front-end

## Kompilera

### Back-end
Bygg hela projektet med:

`make build`

### Front-end

## Testa

### Back-end
Kör alla tester

`make test`

Kör alla tester med verbos flagga, om något test misslyckas skrivs det ut vart det blir fel:

`make testv`

### Front-end

## Starta systemet

### Back-end
För att starta servern kör:

`make run` TODO: Fix this rule

### Front-end

## Struktur

Projektet består av följande kataloger.

### doc

projektrapporter och andra viktiga dokument.

### meta

- Presentation av gruppens medlemmar.
- Gruppkontrakt.
- Projektdagböcker.
- Reflektioner på gruppens arbete.

### api, database, server

All källkod.
