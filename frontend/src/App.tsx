import { useState } from 'react'

import type { Camera } from './types/Camera';
import { CameraModal } from './components/CameraModal';
import { CameraCard } from './components/CameraCard';
import { mockCameras } from './data/mockCameras';

function App() {
  const [selectedCamera, setSelectedCamera] = useState<Camera | null>(null);
  console.log(selectedCamera);

  return (
    <div id="bg" className="h-screen w-screen bg-slate-800">
      <div className="flex h-[6vh] w-full bg-slate-900 items-center justify-around">
        <div className="text-slate-400 w-fit">PLACEHOLDER</div>
        <div className="text-slate-400 w-fit">PLACEHOLDER</div>
        <div className="text-slate-400 w-fit">PLACEHOLDER</div>
      </div>
      <div className="mx-auto h-screen w-[80vw]">
        <div className="p-[.5rem] h-fit bg-slate-700 grid grid-cols-2 xl:grid-cols-3 gap-3">
          {mockCameras.map((camera, i) => {
            return (
                <CameraCard key={camera.id} camera={camera} onClick={() => {setSelectedCamera(camera)}} index={i}></CameraCard>
            )
          })}
        </div>
      </div>

      {selectedCamera &&
        <div className="absolute inset-0 items-center justify-center">
          <CameraModal camera={selectedCamera} onClose={() => setSelectedCamera(null)}></CameraModal>
        </div>
      }
    </div>
  )
}

export default App
