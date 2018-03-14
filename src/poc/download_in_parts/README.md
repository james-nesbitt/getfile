POC: Download in parts

A proof of concept, for golang http client use to retrieve a file in parts.

This POC is heavily based on:
 * https://coderwall.com/p/uz2noa/fast-parallel-downloads-in-golang-with-accept-ranges-and-goroutines

This POC demonstrates:
 * Concurrent part downloads
 * http client & range header

This POC does not:
 * allow partial downloads
 * have any reusable components
 * have any tests (not in focus for the POC)
 *  would have pretty bad memory hits for large files (all downloaded pieces are kept in memory during collection
