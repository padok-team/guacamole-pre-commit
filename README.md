# guacamole-pre-commit

Hook [`pre-commit`](https://pre-commit.com/) pour les fichiers Terraform et
Terragrunt. Le dépôt fournit un exécutable Go qui est destiné à lancer les
contrôles statiques de [Guacamole](https://github.com/padok-team/guacamole)
sur les fichiers concernés.

> **État actuel :** le hook reçoit les chemins des fichiers `.tf` et `.hcl`,
> mais n'exécute pas encore les contrôles Guacamole. Pour l'instant, il affiche
> les fichiers reçus. Il ne bloque donc pas les commits en fonction de leur
> contenu.

## Utilisation dans un projet

Le hook s'ajoute dans le dépôt qui contient le code Terraform/Terragrunt, dans
son fichier `.pre-commit-config.yaml` :

```yaml
repos:
  - repo: https://github.com/padok-team/guacamole-pre-commit
    rev: <version-ou-sha>
    hooks:
      - id: guacamole-pre-commit
```

Remplacez `<version-ou-sha>` par le tag d'une version publiée, ou par un SHA de
commit. Il est préférable d'épingler une version plutôt qu'une branche pour
garder des exécutions reproductibles.

Installez ensuite `pre-commit` si nécessaire, puis activez le hook dans le
dépôt :

```bash
pre-commit install
```

Lors d'un commit, `pre-commit` transmet au hook les fichiers concernés qui
correspondent à `\.(tf|hcl)$`. Pour lancer le hook manuellement sur les fichiers
modifiés :

```bash
pre-commit run guacamole-pre-commit
```

Pour le lancer sur tous les fichiers correspondants du dépôt :

```bash
pre-commit run guacamole-pre-commit --all-files
```

## Développement

Prérequis : Go (version indiquée dans [`go.mod`](./go.mod)).

```bash
make build
make test
make lint
```

`make build` construit le binaire `guacamole-pre-commit` à la racine du dépôt.
`make test` lance les tests Go et `make lint` lance `gofmt` et `golangci-lint`.
