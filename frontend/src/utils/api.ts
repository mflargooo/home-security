import type { Camera } from "../types/Camera"

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

export async function createCamera(data : Camera) {
  const res = await fetch(`${API_URL}/cameras`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error(`Failed to create camera: ${res.statusText}`);
  return res.json();
}

export async function updateCamera(id : string, data : Camera) {
  const res = await fetch(`${API_URL}/cameras/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error(`Failed to update camera: ${res.statusText}`);
  return res.json();
}

export async function deleteCamera(id : string) {
  const res = await fetch(`${API_URL}/cameras/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error(`Failed to delete camera: ${res.statusText}`);
}

// Stream URL helper - converts RTSP to HLS proxy endpoint
// Your backend should expose HLS streams at /api/cameras/:id/stream
export function getStreamUrl(cameraId : string) {
  return `${HLS_URL}/${cameraId}/index.m3u8`;
}
