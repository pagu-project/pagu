# Debugging

## Debugging Engine

Pagu is designed to keep the internal logic independent from the user interface layer.
Therefore, the best and simplest way to debug the engine and command functionalities is by using the **CLI (Command Line Interface)**.
This method gives you full access to all commands and the engine internals.

## Debugging User Interface

### Discord Platform

Run Pagu with the required Discord configuration and start local debugging.

### Telegram Platform

Run Pagu with the required Telegram configuration and start local debugging.

### WhatsApp Platform

Debugging Pagu on WhatsApp is a bit more complex. WhatsApp requires a callback URL that supports SSL.
The simplest approach is to set up a webhook like `https://my-website.com/webhook` and
use **Nginx reverse proxy** combined with **SSH tunneling** to redirect traffic to a local machine.

The Nginx block configuration might look like this:

```nginx
location /webhook {
    proxy_pass http://127.0.0.1:3000/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

Then, establish an SSH tunnel with:

```bash
ssh -R 3000:localhost:3000 user@my-website.com
```
