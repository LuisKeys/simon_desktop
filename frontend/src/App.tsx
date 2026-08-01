import {useCallback, useEffect, useState} from 'react'
import './App.css'
import {CreateConversation, ListConversations} from '../wailsjs/go/app/AppService'
import {conversations} from '../wailsjs/go/models'
import Chat from './Chat'
import SettingsMenu from './settings/SettingsMenu'

function formatUpdatedAt(value: string): string {
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) {
        return ''
    }
    return date.toLocaleString(undefined, {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}

function App() {
    const [conversationList, setConversationList] = useState<conversations.Conversation[]>([])
    const [activeId, setActiveId] = useState<string | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    const refreshConversations = useCallback(() => {
        ListConversations()
            .then((list) => setConversationList(list ?? []))
            .catch(() => setError('Could not load conversations.'))
    }, [])

    useEffect(() => {
        ListConversations()
            .then((list) => {
                setConversationList(list ?? [])
                if (list && list.length > 0) {
                    setActiveId(list[0].id)
                }
            })
            .catch(() => setError('Could not load conversations.'))
            .finally(() => setLoading(false))
    }, [])

    async function handleNewChat() {
        setError(null)
        try {
            const created = await CreateConversation()
            setConversationList((prev) => [created, ...prev])
            setActiveId(created.id)
        } catch {
            setError('Could not create a new conversation.')
        }
    }

    const activeConversation = conversationList.find((c) => c.id === activeId) ?? null

    return (
        <div id="App">
            <div id="skeleton">
                <div id="sidebar">
                    <div id="logo-mark">
                        <div className="logo-mark-brand" aria-hidden="true">
                            <svg viewBox="0 0 100 100">
                                <polygon
                                    points="50,4 90,27 90,73 50,96 10,73 10,27"
                                    fill="none"
                                    stroke="#FF6A00"
                                    strokeWidth="4"
                                />
                                <path
                                    d="M32 40 C32 25, 40 20, 40 32 C44 24, 56 24, 60 32 C60 20, 68 25, 68 40 Z"
                                    fill="none"
                                    stroke="var(--color-text)"
                                    strokeWidth="3"
                                />
                                <circle cx="40" cy="46" r="2.5" fill="var(--color-text)"/>
                                <circle cx="60" cy="46" r="2.5" fill="var(--color-text)"/>
                                <text x="50" y="62" textAnchor="middle" fill="#FF6A00" fontSize="14" fontFamily="monospace">&gt;_</text>
                            </svg>
                            <span>Simon</span>
                        </div>
                        <SettingsMenu/>
                    </div>

                    <button id="new-chat" onClick={handleNewChat}>+ New chat</button>

                    <div id="conversation-list">
                        {loading && <div className="muted">Loading…</div>}
                        {!loading && conversationList.length === 0 && (
                            <div className="muted">No conversations yet</div>
                        )}
                        {conversationList.map((c) => (
                            <button
                                key={c.id}
                                className={'conversation-item' + (c.id === activeId ? ' active' : '')}
                                onClick={() => setActiveId(c.id)}
                            >
                                <span className="conversation-title">{c.title}</span>
                                <span className="conversation-date">{formatUpdatedAt(c.updatedAt as unknown as string)}</span>
                            </button>
                        ))}
                    </div>

                    <div id="sidebar-footer">
                        <button id="settings-button" disabled title="Coming in a later phase">
                            Documents
                        </button>
                    </div>
                </div>

                <div id="main">
                    {error && <div className="error-banner">{error}</div>}
                    {activeConversation ? (
                        <Chat
                            key={activeConversation.id}
                            conversationId={activeConversation.id}
                            conversationTitle={activeConversation.title}
                            onActivity={refreshConversations}
                        />
                    ) : (
                        <div id="chat-placeholder">
                            <p>No conversation selected.</p>
                            <p className="muted">Start one with "+ New chat".</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

export default App
