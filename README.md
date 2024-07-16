# Battlesship

## Connect

### Locally

```bash
ssh -p 23234 localhost
```

### Remotely

```bash
ssh battleuser@battle.orb.local -p 23234
```

### Build the guy

1. Create the image

```bash
❯ docker build --tag 'battlesship-image' .
```

2. Run the container

(TODO: should add the detach later...)

```bash
❯ docker run -p 23234:23234 --name battle 'battlesship-image'
```
