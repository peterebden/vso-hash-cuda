// This is a very minimal main designed to run a pure cuda (i.e. no Go) kernel.
// It's only useful for troubleshooting device setup etc.

#include <stdio.h>
#include <string.h>

#include "sha256.h"

static const int num_jobs = 4;
static const char* inputs[] = {
  "The most merciful thing in the world, I think, is the inability of the human mind to correlate all its contents.",
  "We live on a placid island of ignorance in the midst of black seas of infinity, and it was not meant that we should voyage far.",
  "The sciences, each straining in its own direction, have hitherto harmed us little; but some day the piecing together of dissociated knowledge will open up such terrifying vistas of reality, and of our frightful position therein, that we shall either go mad from the revelation or flee from the deadly light into the peace and safety of a new dark age.",
  "Theosophists have guessed at the awesome grandeur of the cosmic cycle wherein our world and human race form transient incidents.",
};

#define CHECK_ERROR() { const char* err = sha256_last_error(); if (err != NULL) { fprintf(stderr, "%s\n", err); return 1; } }

int main(int argc, char **argv) {
  sha256_preinit();
  CHECK_ERROR();
  SHA256_job** jobs = sha256_alloc_jobs(num_jobs, 2048);  // 2048 is arbitrarily bigger than any of the sentences above.
  if (!jobs) {
    CHECK_ERROR();
  }
  for (int i = 0; i < num_jobs; ++i) {
    if (sha256_init_job(jobs, i, (const unsigned char*)inputs[i], strlen(inputs[i]))) {
      CHECK_ERROR();
    }
  }

  sha256_run(jobs, num_jobs);
  CHECK_ERROR();

  for (int i = 0; i < num_jobs; ++i) {
    printf("Sentence %d: ", i);
    for (int j = 0; j < 32; ++j) {
      printf("%.2x", jobs[i]->digest[j]);
    }
    printf("\n");
  }

  sha256_free_jobs(jobs, num_jobs);
  CHECK_ERROR();
  return 0;
}
