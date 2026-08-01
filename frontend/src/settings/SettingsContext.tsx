import {createContext, useContext, useEffect, useState, ReactNode} from 'react'
import {GetSetting, SetSetting, SetWindowScale} from '../../wailsjs/go/app/AppService'

export type ThemeMode = 'dark' | 'light'

const THEME_KEY = 'ui.theme'
const SCALE_KEY = 'ui.scale'
const DEFAULT_THEME: ThemeMode = 'dark'
const DEFAULT_SCALE = 1

interface SettingsContextValue {
    theme: ThemeMode
    scale: number
    setTheme: (t: ThemeMode) => void
    setScale: (s: number) => void
}

const SettingsContext = createContext<SettingsContextValue | null>(null)

export function SettingsProvider({children}: {children: ReactNode}) {
    const [theme, setThemeState] = useState<ThemeMode>(DEFAULT_THEME)
    const [scale, setScaleState] = useState<number>(DEFAULT_SCALE)

    useEffect(() => {
        Promise.all([GetSetting(THEME_KEY), GetSetting(SCALE_KEY)])
            .then(([t, s]) => {
                if (t === 'light' || t === 'dark') setThemeState(t)
                const parsed = parseFloat(s)
                if (!Number.isNaN(parsed) && parsed > 0) setScaleState(parsed)
                SetWindowScale(!Number.isNaN(parsed) && parsed > 0 ? parsed : DEFAULT_SCALE).catch(() => {})
            })
            .catch(() => {})
    }, [])

    useEffect(() => {
        document.documentElement.dataset.theme = theme
    }, [theme])

    useEffect(() => {
        document.documentElement.style.setProperty('--ui-scale', String(scale))
    }, [scale])

    function setTheme(t: ThemeMode) {
        setThemeState(t)
        SetSetting(THEME_KEY, t).catch(() => {})
    }

    function setScale(s: number) {
        setScaleState(s)
        SetSetting(SCALE_KEY, String(s)).catch(() => {})
        SetWindowScale(s).catch(() => {})
    }

    return (
        <SettingsContext.Provider value={{theme, scale, setTheme, setScale}}>
            {children}
        </SettingsContext.Provider>
    )
}

export function useSettings(): SettingsContextValue {
    const ctx = useContext(SettingsContext)
    if (!ctx) {
        throw new Error('useSettings must be used within SettingsProvider')
    }
    return ctx
}
