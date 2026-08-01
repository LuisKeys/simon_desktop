import {useEffect, useRef, useState} from 'react'
import {useSettings} from './SettingsContext'
import {Quit} from '../../wailsjs/runtime/runtime'
import './SettingsMenu.css'

const SCALE_MIN = 1
const SCALE_MAX = 2
const SCALE_STEP = 0.1

export default function SettingsMenu() {
    const {theme, scale, setTheme, setScale} = useSettings()
    const [open, setOpen] = useState(false)
    const containerRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        function onClickOutside(e: MouseEvent) {
            if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
                setOpen(false)
            }
        }
        document.addEventListener('mousedown', onClickOutside)
        return () => document.removeEventListener('mousedown', onClickOutside)
    }, [])

    return (
        <div id="settings-menu" ref={containerRef}>
            <button
                id="settings-gear"
                aria-label="Settings"
                aria-expanded={open}
                onClick={() => setOpen((v) => !v)}
            >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
                </svg>
            </button>

            {open && (
                <div id="settings-dropdown">
                    <div className="settings-row">
                        <span>Theme</span>
                        <button
                            className="theme-toggle"
                            role="switch"
                            aria-checked={theme === 'light'}
                            onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                        >
                            <span className="theme-toggle-thumb"/>
                        </button>
                    </div>
                    <div className="settings-row">
                        <span>Scale</span>
                        <input
                            type="range"
                            min={SCALE_MIN}
                            max={SCALE_MAX}
                            step={SCALE_STEP}
                            value={scale}
                            onChange={(e) => setScale(parseFloat(e.target.value))}
                        />
                        <span className="settings-scale-value">{Math.round(scale * 100)}%</span>
                    </div>
                    <button className="settings-quit" onClick={() => Quit()}>
                        Quit
                    </button>
                </div>
            )}
        </div>
    )
}
