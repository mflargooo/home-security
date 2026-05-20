import { EllipsisVertical, X } from "lucide-react"

type ModalBarProps = {
    title: string;
    xOnClick: (() => void) | null;
    ellipsisOnClick: (() => void) | null;
}

export function ModalBar({ title, xOnClick, ellipsisOnClick } : ModalBarProps) {
    return (
        <div className="flex text-slate-400 p-[.375rem] lg:p-[.5rem] justify-between">
            <button className={`${xOnClick ? "" : "opacity-0 pointer-events-none"} bg-red-700 border-red-900 border rounded-[100%] hover:bg-red-500 w-fit h-fi my-auto`} onClick={() => xOnClick()}><X/>
            </button>
            <span className="font-bold text-lg md:text-xl lg:text-2xl"> {title} </span>
            <button className={`${ellipsisOnClick ? "" : "opacity-0 pointer-events-none"} hover:bg-slate-500/50 border rounded-[100%] border-none w-fit h-fit my-auto`}><EllipsisVertical/>
            </button>
        </div>
    )
}