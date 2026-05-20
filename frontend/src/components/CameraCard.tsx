import type { Camera } from "../types/Camera"

type CameraCardProps = {
    camera: Camera;
    onClick: (camera: Camera) => void;
    index: number;
}
export function CameraCard({ camera, onClick, index } : CameraCardProps) {
    const { id, name, status, metadata, last_seen_at } = camera;
    index;
    return (
        <div 
            className="
                group relative aspect-video bg-slate-800 border border-slate-900 border-3 rounded-md text-red-500
                hover:bg-slate-700 transition-all duration-200 overflow-hidden p-[.5rem] hover:border-accent 
            " 
        >
            <div className="absolute inset-0 p-[1rem]">
                <div className="justify-between">
                    <div className="font-bold text-slate-500 text-md xl:text-lg"> {camera.name} </div>
                    <div className="italic text-slate-500 text-xs xl:text-sm mt-[-.125rem]"> {camera.id.toUpperCase()} </div>
                </div>
            </div>
            <div 
                className="opacity-0 group-hover:opacity-100 transition-opacity absolute flex inset-0 items-center justify-center"
                onClick={() => onClick(camera)}
            >
                <span className="font-bold text-lg xl:text-2xl text-slate-400">View Stream</span>
            </div>
        </div>
    )
}