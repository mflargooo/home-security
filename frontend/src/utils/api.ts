import type { Camera } from "../types/Camera"
import type { Clip } from "../types/Clip"

// Replace API_URL with your actual API endpoint
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/';
const HLS_URL = import.meta.env.VITE_HLS_URL || 'http://localhost:8888/';

export async function fetchCameras() {
  const res = await fetch(`${API_URL}/cameras`);
  if (!res.ok) throw new Error(`Failed to fetch cameras: ${res.statusText}`);
  return res.json();
}

export async function fetchCamera(id : string) {
  const res = await fetch(`${API_URL}/cameras/${id}`);
  if (!res.ok) throw new Error(`Failed to fetch camera: ${res.statusText}`);
  return res.json();
}

export async function createCamera(camera : Camera) {
  const res = await fetch(`${API_URL}/cameras`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(camera),
  });
  if (!res.ok) throw new Error(`Failed to create camera: ${res.statusText}`);
  return res.json();
}

export async function updateCamera(camera : Camera) {
  const res = await fetch(`${API_URL}/cameras/${camera.id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(camera),
  });
  if (!res.ok) throw new Error(`Failed to update camera: ${res.statusText}`);
  return res.json();
}

export async function deleteCamera(id : string) {
  const res = await fetch(`${API_URL}/cameras/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error(`Failed to delete camera: ${res.statusText}`);
}

export async function createSnapshot(camera: Camera) {
  const res = await fetch(`${API_URL}/snapshots/${camera.id}`, {
    method: 'POST'
  })
  if (!res.ok) throw new Error(`Failed to create snapshot: ${res.statusText}`);
  return res.json()
}

export async function deleteSnapshot(sessionID: string) {
  const res = await fetch(`${API_URL}/snapshots/${sessionID}`, {
    method: 'DELETE'
  })
  if (!res.ok) throw new Error(`Failed to delete snapshot: ${res.statusText}`);
  return res
}

export async function createClip(sessionID: string, clip: Clip) {
  const res = await fetch(`${API_URL}/clips/${sessionID}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(clip),
  })
  if (!res.ok) throw new Error(`Failed to create clip: ${res.statusText}`);
  return res.json()
}

// Stream URL helper - converts RTSP to HLS proxy endpoint
// Your backend should expose HLS streams at /api/cameras/:id/stream
export function getStreamUrl(cameraID : string) {
  return `${HLS_URL}/${cameraID}/index.m3u8`;
}

export function getSnapshotUrl(uri : string) {
  return `${API_URL}/${uri}`;
}
