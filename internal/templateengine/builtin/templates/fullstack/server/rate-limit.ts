type Bucket = { count: number; resetsAt: number };
const buckets = new Map<string, Bucket>();

export function allowRequest(key: string, limit = 60, windowMs = 60_000) {
  const now = Date.now();
  const bucket = buckets.get(key);
  if (!bucket || bucket.resetsAt <= now) {
    buckets.set(key, { count: 1, resetsAt: now + windowMs });
    return true;
  }
  bucket.count += 1;
  return bucket.count <= limit;
}

// This limiter is correct for one process only. Replace it with a shared,
// atomic store before horizontally scaling or protecting a high-risk route.
