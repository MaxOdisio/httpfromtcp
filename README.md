# httpfromtcp

Un server **HTTP/1.1 scritto da zero in Go**, partendo da un socket TCP grezzo, senza usare `net/http`.

![Go](https://img.shields.io/badge/Go-1.25.6-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)

Il progetto è nato per capire cosa succede *sotto* un framework web: come si leggono i byte da una connessione TCP, come si ricostruisce una richiesta HTTP da un flusso, e come si costruisce una risposta conforme alle specifiche (RFC 9110).

> Progetto didattico, sviluppato seguendo il percorso "From TCP to HTTP" di Boot.dev e poi rifinito per conto mio. Non è pensato per la produzione (vedi [Limiti](#limiti)).

---

## Cosa implementa

- **Lettura a stream**: i dati vengono processati man mano che arrivano sul buffer, senza aspettare l'intera richiesta.
- **Server TCP concorrente**: ciclo `Accept` con una goroutine per ogni connessione.
- **Parsing manuale** di request line (metodo, URI, versione), header e body.
- **Header**: tipo `Headers` dedicato, con chiavi normalizzate in lowercase all'inserimento.
- **Risposte**: scrittura di status line, header e body tramite un `Writer` con macchina a stati.
- **Chunked Transfer Encoding** per lo streaming delle risposte.
- **Dati binari** nel body.

## Struttura del progetto

```
.
├── cmd/            # entrypoint eseguibili
└── internal/
    ├── request/    # parsing della richiesta (state machine)
    ├── headers/    # tipo Headers e normalizzazione delle chiavi
    ├── response/   # Writer per status line, header e body
    └── server/     # listener TCP e gestione delle connessioni
```

## Come eseguirlo

Requisiti: Go installato (vedi `go.mod` per la versione).

```bash
git clone https://github.com/MaxOdisio/httpfromtcp.git
cd httpfromtcp

# avvia il server
go run ./cmd/httpserver
```

In un altro terminale:

```bash
# richiesta semplice
curl -i http://localhost:42069/

# risposta in chunked encoding
curl -i --raw http://localhost:42069/httpbin/stream/5
```

## Test

```bash
go test ./...
go test -race ./...
```

## Scelte di design e cose imparate

- **Parsing incrementale con state machine**: la richiesta può arrivare spezzata in più `Read`, quindi il parser tiene lo stato e consuma solo i byte completi.
- **Data race sugli header**: una mappa condivisa tra goroutine causava una race; risolta evitando la mutazione condivisa.
- **Normalizzazione delle chiavi**: gli header HTTP sono case-insensitive, quindi le chiavi vengono portate in lowercase *all'inserimento* e non al momento della lettura, così `Set`, `Get` e `Replace` restano coerenti.
- **Controllo degli errori**: ogni scrittura sul `Writer` restituisce un errore che viene sempre controllato, perché una connessione può cadere in qualsiasi momento.

## Limiti

Cose che il progetto **non** copre (di proposito):

- niente keep-alive / connessioni persistenti
- niente HTTP/2 né TLS
- nessun limite sulla dimensione di header e body
- copertura parziale di RFC 7230

## Possibili sviluppi

- [ ] Supporto a keep-alive
- [ ] Timeout di lettura/scrittura sulle connessioni
- [ ] Benchmark contro `net/http`

## Licenza

Distribuito con licenza MIT. Vedi il file [LICENSE](LICENSE).
