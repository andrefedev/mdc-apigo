# plan.md

## Objetivo

Consolidar la arquitectura actual del monolito sin volver a una separacion por feature que ya no aporta al proyecto.

## Direccion elegida

- un solo `ApiService`
- un solo nucleo en `internal/app`
- un borde gRPC en `api/okgrpc`
- bootstrap explicito en `cmd/server/main.go`

## Prioridad P0

- mantener documentacion y codigo alineados con `internal/app`
- evitar que nuevos cambios revivan `internal/features/*` como patron
- mantener auth y errores del transporte centralizados en `api/okgrpc`

## Prioridad P1

- terminar de mover cualquier resto legacy fuera de estructuras viejas
- estabilizar convenciones de `datax -> data -> repository`
- mantener un unico criterio de autorizacion por metodo

## Prioridad P2

- estabilizar nombres del `ApiService`
- reducir drift entre schema real y modelos de `internal/app`
- agregar tests focalizados del borde gRPC y de `UseService`

## Criterio de exito

- el borde vive en `api/*`
- el nucleo vive en `internal/app`
- el contrato externo sigue siendo un solo `ApiService`
- no hay ambiguedad sobre donde va un cambio nuevo
