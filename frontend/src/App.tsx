import { useState, useEffect } from 'react'

import type { Camera } from './types/Camera';
import { CameraStreamModal } from './components/CameraStreamModal';
import { CameraCard } from './components/CameraCard';
import { fetchCameras } from './utils/api';

function App() {
  const [selectedCamera, setSelectedCamera] = useState<Camera | null>(null);
  const [cameras, setCameras] = useState<Camera[]>([]);

  const loadCameras = async () => {
    const data = await fetchCameras();
    setCameras(data);
  }

  useEffect(() => {
    loadCameras()
    return () => {}
  }, [])

  return (
    <div id="bg" className="h-screen w-screen bg-slate-800">
      <div className="flex h-[6vh] w-full bg-slate-900 items-center justify-around">
        <div className="text-slate-400 w-fit">PLACEHOLDER</div>
        <div className="text-slate-400 w-fit">PLACEHOLDER</div>
        <div className="text-slate-400 w-fit">PLACEHOLDER</div>
      </div>
      <div className="mx-auto h-screen w-[80vw]">
        <div className="p-[.5rem] h-fit bg-slate-700 grid grid-cols-2 xl:grid-cols-3 gap-3">
          {cameras.map((camera, i) => {
            if (camera.name === '') { // set default camera name
              camera.name = "CAMERA-" + i.toString().padStart(3, "0")
            };
            return (
                <CameraCard key={camera.id} camera={camera} onClick={() => {setSelectedCamera(camera)}} index={i}></CameraCard>
            )
          })}
        </div>
      </div>

      {selectedCamera &&
        <div className="absolute inset-0 items-center justify-center">
          <CameraStreamModal camera={selectedCamera} onClose={() => setSelectedCamera(null)}></CameraStreamModal>
        </div>
      }
    </div>
  )
}

export default App
