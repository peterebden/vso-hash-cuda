#ifndef SHA256_H
#define SHA256_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

typedef unsigned char* SHA256_data;

typedef struct SHA256_job {
	unsigned char* data;
	unsigned long long size;
	unsigned char digest[64];
} SHA256_job;

// Runs some initialisation. Must be called once before sha256_run is called.
void sha256_preinit();

// Calculates SHA256 against the given inputs.
void sha256_run(SHA256_job** jobs, int n);

// Allocates a number of jobs, each of the given size.
// Returns NULL on failure.
SHA256_job** sha256_alloc_jobs(int n, int size);

// Frees previously allocated jobs.
void sha256_free_jobs(SHA256_job** jobs, int n);

// Initialises the i'th job with the given data.
// The size of data must be <= the size given to sha256_alloc_jobs.
// Returns 0 on success.
int sha256_init_job(SHA256_job** jobs, int i, const unsigned char* data, int size);

// Returns a description of the last error occurring, or NULL if no error occurred.
// The caller should not free the buffer.
const char* sha256_last_error();

// Copies the digest from a job into the given array. Convenience method to help cgo.
void sha256_copy_digest(SHA256_job** jobs, int i, unsigned char* dest);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // SHA256_H
