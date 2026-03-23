// Package apperr define el error canónico de aplicación usado por servicios,
// repositorios y adaptadores de transporte.
//
// El contrato actual del paquete separa dos capas:
//   - metadatos técnicos para propagación y diagnóstico: Op, Kind y Cause
//   - contrato público opcional para cliente: Code y Body
//
// Uso esperado por capa:
//   - repository: retorna errores técnicos con Internal, NotFound, Conflict, etc.
//   - service: envuelve con Wrap y, cuando aplica, agrega contrato público con WithPublic
//   - transporte: convierte el error con ResponseOf y KindOf
//
// ResponseOf garantiza un payload público estable incluso cuando el error no
// define Code o Body explícitos, usando defaults derivados de Kind.
package apperr
