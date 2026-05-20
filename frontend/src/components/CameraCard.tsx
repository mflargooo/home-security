import type { Camera } from "../types/Camera"

type CameraCardProps = {
    camera: Camera;
    onClick: (camera: Camera) => void;
    index: number;
}
export function CameraCard({ camera, onClick, index } : CameraCardProps) {
    index;
    return (
        <div 
            className="
                group relative aspect-video bg-slate-800 border border-slate-900 border-3 rounded-md text-red-500
                transition-all duration-200 overflow-hidden p-[.25rem] hover:border-accent 
            " 
        >
            <div id="thumbnail" className="bg-slate-800 w-full h-full rounded-sm">

            </div>
            <div className="absolute bottom-0 left-0 p-[.5rem] w-full bg-slate-900/90">
                <div className="justify-between">
                    <div className="font-bold text-slate-500 text-base"> {camera.name} </div>
                    <div className="italic text-slate-500 text-xs mt-[-.125rem]"> {camera.id.toUpperCase()} </div>
                </div>
            </div>
            <div 
                className="absolute flex inset-0 items-center justify-center opacity-0 group-hover:opacity-100 group-hover:bg-slate-700/30  transition-opacity"
                onClick={() => onClick(camera)}
            >
                <span className="font-bold text-lg xl:text-2xl text-slate-400">View Stream</span>
            </div>
        </div>
    )
}