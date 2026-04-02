# plan.md

## Objetivo

Consolidar la arquitectura actual sin volver a dispersar el transporte dentro de las features.

## Direccion elegida

- un solo `ApiService`
- `api/*` como borde canonico
- `internal/features/*` como nucleo por feature
- interceptors y auth centralizados en el transporte

## Prioridad P0

- mantener documentacion y codigo alineados
- evitar nuevos handlers gRPC dentro de `internal/features/*`
- evitar dividir protobuf en multiples services sin necesidad

## Prioridad P1

- terminar de mover restos de transporte legacy fuera de features
- dejar auth del transporte en una sola estrategia
- mantener el mapping gRPC en `api/okgrpc`

## Prioridad P2

- estabilizar nombres y consistencia del `ApiService`
- reducir drift entre schema real y features
- agregar tests focalizados del borde gRPC

## Criterio de exito

- el borde vive en `api/*`
- el dominio vive en `internal/features/*`
- el contrato externo sigue siendo un solo `ApiService`
- no hay ambiguedad sobre donde va cada cambio nuevo
