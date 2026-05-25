# HTTP/1.1 Server from Scratch 🚀

Questo repository contiene la mia implementazione da zero del protocollo HTTP/1.1 scritta in Go.

L'obiettivo di questo progetto non è usare il pacchetto `net/http` standard, ma capire come muovere, parsare e strutturare i flussi di byte grezzi partendo da un socket TCP di basso livello, fino ad arrivare a gestire le specifiche RFC 7230.

> ⚠️ **Project Status: Work in Progress** > Sto attivamente sviluppando il server.

---

## 🧠 Cosa sto imparando (e implementando)

- **HTTP Streams:** Come leggere e processare flussi di byte (streaming) man mano che arrivano sul buffer senza bloccare il server.
- **TCP Layer:** Capire la differenza profonda tra TCP e UDP e gestire le connessioni concorrenti.
- **Parsing manuale:** Analizzare a basso livello la *Request Line* (Metodo, URI, Versione) e mappare gli *HTTP Headers* senza framework esterni.
- **Gestione dei Body e Encoding:** Leggere il payload delle richieste e implementare il supporto al *Chunked Transfer Encoding*.

---

## 🛠️ Roadmap delle Funzionalità

Di seguito trovi lo stato di avanzamento del progetto:

| Capitolo | Funzionalità / Topic | Stato | Note |
| :---: | :--- | :---: | :--- |
| **1** | **HTTP Streams & Byte Chunking** | 🟢 Completato | Lettura controllata a blocchi e gestione dei `\n`. |
| **2** | **TCP Sockets** | 🟢 Completato | Gestione del ciclo di `Accept` e connessioni TCP affidabili. |
| **3** | **HTTP Requests Intro** | 🟢 Completato | Studio della struttura anatomica di una richiesta web. |
| **4** | **Request Lines Parsing** | 🟢 Completato | Estrazione di Metodo (GET/POST), URI e Versione. |
| **5** | **HTTP Headers** | 🟡 In Corso | Mappatura chiave-valore degli header. |
| **6** | **HTTP Body** | 🔴 Pianificato | Lettura e processing del payload della richiesta. |
| **7** | **HTTP Responses** | 🔴 Pianificato | Generazione e invio di risposte conformi al client. |
| **8** | **Chunked Encoding** | 🔴 Pianificato | Streaming dei dati a segmenti. |
| **9** | **Binary Data** | 🔴 Pianificato | Gestione dei dati binari e versioni del protocollo. |

---

## 💻 Tech Stack & Requisiti

- **Linguaggio:** Go (Golang)
- **Moduli:** Standard Library (`net` per i socket TCP, `io`, `bufio`, `bytes`)
- **Testing suite:** Boot.dev CLI tool per la validazione locale delle challenge

---

## 🏃 Configurazione Locale ed Esecuzione

Se vuoi dare un'occhiata a come gestisco lo streaming dei byte:

1. Assicurati di avere Go installato sul tuo sistema.
2. Clona il repository:
   ```bash
   git clone https://github.com/MaxOdisio/httpfromtcp.git
   cd httpfromtcp
