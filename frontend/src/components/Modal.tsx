    import type { ReactNode } from "react";
    import { EllipsisVertical, X } from "lucide-react"

    type ModalBarProps = {
        title: string;
        onClose: (() => void);
        ellipsisOnClick?: (() => void);
    }

    type ModalProps = ModalBarProps & {
        hidden?: boolean
        className?: string
        children?: ReactNode;
    }

    export function ModalBar({ title, onClose, ellipsisOnClick } : ModalBarProps) {
        return (
            <div className="flex bg-slate-900 text-slate-400 p-[.375rem] lg:p-[.5rem] justify-between">
                <button className={`bg-red-700 border-red-900 border rounded-[100%] hover:bg-red-500 w-fit h-fi my-auto`} onClick={() => onClose()}><X/>
                </button>
                <span className="font-bold text-md md:text-lg lg:text-xl"> {title} </span>
                <button className={`${ellipsisOnClick ? "" : "opacity-0 pointer-events-none"} hover:bg-slate-500/50 border rounded-[100%] border-none w-fit h-fit my-auto`} onClick={ellipsisOnClick ? () => ellipsisOnClick() : () => {}}><EllipsisVertical/>
                </button>
            </div>
        )
    }

    export function Modal({ title, onClose, ellipsisOnClick, hidden=false, className, children } : ModalProps ) {
        
        return (
            <div hidden={hidden}
                className="fixed inset-0 z-50 flex modal-backdrop items-center justify-center bg-slate-900/80" 
                onClick={(e) => e.target === e.currentTarget ? onClose() : (() => {})()}
            >
                <div 
                    className={`
                        relative inset-0 m-auto bg-slate-800 
                        border border-rounded rounded-xl overflow-hidden
                        shadow-2xl h-fit ${className}
                    `}
                >
                <ModalBar title={title} onClose={() => onClose()} ellipsisOnClick={ellipsisOnClick} />
                { children }
                </div>
            </div>
        )
    }