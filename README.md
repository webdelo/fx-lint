# fx-lint

`fx-lint` — это набор [`go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis) аналайзеров, которые обеспечивают соблюдение конвенций оформления [`uber-go/fx`](https://github.com/uber-go/fx)-обвязки, принятых в проектах [webdelo](https://github.com/webdelo).

## Зачем этот пакет

Конвенции UFX (UFX-002, UFX-024, UFX-026) определяют, **где** и **как** в проекте может появляться код, использующий `go.uber.org/fx`:

- вся wiring-логика (`fx.Module`, `fx.Provide`, `fx.Invoke`) должна жить только под `internal/infrastructure/fx/modules/`, `internal/app/` и `internal/cmd/`;
- orchestration-файлы внутри `fxm_*`-пакетов (`module.go`, `providers.go`, `lifecycle.go`) обязаны содержать только одну одноимённую функцию-сборщик; конструкторы — в `fxs_*`/`fxh_*`;
- в `fx.Provide`/`fx.Invoke` нельзя передавать inline-лямбды — только именованные функции, чтобы fx-граф зависимостей оставался читаемым в debug-выводе.

`fx-lint` автоматизирует эти проверки и распространяется как самостоятельный Go-модуль, который можно:

- запустить как standalone-multichecker (`bin/fxlint`),
- встроить в [`golangci-lint custom`](https://golangci-lint.run/plugins/module-plugins/) как Module Plugin (v2).

## Правила

| Имя в golangci-lint | Конвенция | Что ловит |
|----|----|----|
| `fxinlinelambda` | UFX-026 | Inline-лямбды в `fx.Provide` / `fx.Invoke` |
| `fxmodulelocation` | UFX-024 | Импорт `go.uber.org/fx` или вызов `fx.Module/Provide/Invoke` из пакета вне whitelist-дерева |
| `fxmoduleorchestration` | UFX-002 / UFX-024 | Конструкторы и методы в orchestration-файлах `fxm_*/module.go`, `providers.go`, `lifecycle.go` |

## Использование

### Standalone-multichecker

```bash
go install github.com/webdelo/fx-lint/cmd/fxlint@latest

fxlint -mod-path=github.com/foo/bar ./...
```

`-mod-path` — модульный путь линтуемого проекта; на его основе вычисляются whitelist-префиксы для `fxmodulelocation`. Также читается из переменной окружения `FX_LINT_MOD_PATH`.

### golangci-lint Module Plugin (v2)

В корне проекта:

**`.custom-gcl.yml`**
```yaml
version: v2.1.0
name: custom-gcl
destination: ./bin
plugins:
  - module: github.com/webdelo/fx-lint
    version: v0.1.0
```

**`.golangci.yaml`** (фрагмент)
```yaml
linters:
  enable:
    - fxinlinelambda
    - fxmodulelocation
    - fxmoduleorchestration
  settings:
    custom:
      fxlint:
        type: module
        description: webdelo fx wiring conventions
        settings:
          mod_path: github.com/foo/bar
```

Собрать кастомный бинарь один раз и использовать его:
```bash
golangci-lint custom              # → ./bin/custom-gcl
./bin/custom-gcl run ./...
```

## Происхождение

Аналайзеры исторически появились в [`github.com/Issengaard/crm_printing_house`](https://github.com/Issengaard/crm_printing_house) (commit `4927720`) и затем были портированы в виде генератора `tools/fxlint` в [`github.com/webdelo/scratch`](https://github.com/webdelo/scratch) (REQ-BUG-FIX-TASK-013). Этот репозиторий — финальная итерация: вынос правил в самостоятельный, версионируемый Go-модуль.

## License

MIT
