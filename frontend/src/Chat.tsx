import {useEffect, useRef, useState} from 'react'
import {CancelRun, GetConversationMessages, SendMessage} from '../wailsjs/go/app/AppService'
import {EventsOn} from '../wailsjs/runtime/runtime'
import {sessions} from '../wailsjs/go/models'

type ChatRunPayload = { conversationId: string }
type ChatCompletedPayload = { conversationId: string; text: string }
type ChatFailedPayload = { conversationId: string; message: string }

interface ChatProps {
    conversationId: string
    conversationTitle: string
    onActivity: () => void
}

function Chat({conversationId, conversationTitle, onActivity}: ChatProps) {
    const [messages, setMessages] = useState<sessions.ChatMessage[]>([])
    const [input, setInput] = useState('')
    const [running, setRunning] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [loadingHistory, setLoadingHistory] = useState(true)
    const bottomRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        setLoadingHistory(true)
        setError(null)
        GetConversationMessages(conversationId)
            .then((list) => setMessages(list ?? []))
            .catch(() => setError('Could not open the conversation.'))
            .finally(() => setLoadingHistory(false))
    }, [conversationId])

    useEffect(() => {
        const offStarted = EventsOn('chat.run.started', (payload: ChatRunPayload) => {
            if (payload.conversationId === conversationId) setRunning(true)
        })
        const offCompleted = EventsOn('chat.run.completed', (payload: ChatCompletedPayload) => {
            if (payload.conversationId !== conversationId) return
            setRunning(false)
            setMessages((prev) => [...prev, {role: 'assistant', content: payload.text}])
            onActivity()
        })
        const offFailed = EventsOn('chat.run.failed', (payload: ChatFailedPayload) => {
            if (payload.conversationId !== conversationId) return
            setRunning(false)
            setError(payload.message)
        })
        const offCancelled = EventsOn('chat.run.cancelled', (payload: ChatRunPayload) => {
            if (payload.conversationId !== conversationId) return
            setRunning(false)
            setError('The request was cancelled.')
        })
        return () => {
            offStarted()
            offCompleted()
            offFailed()
            offCancelled()
        }
    }, [conversationId, onActivity])

    useEffect(() => {
        bottomRef.current?.scrollIntoView({behavior: 'smooth'})
    }, [messages, running])

    async function handleSend() {
        const text = input.trim()
        if (!text || running) return
        setInput('')
        setError(null)
        setMessages((prev) => [...prev, {role: 'user', content: text}])
        try {
            await SendMessage(conversationId, text)
        } catch (err) {
            setRunning(false)
            setError(typeof err === 'string' ? err : 'Could not send the message.')
        }
    }

    function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault()
            void handleSend()
        }
    }

    async function handleCancel() {
        try {
            await CancelRun(conversationId)
        } catch {
            // best-effort: the run may have already finished
        }
    }

    return (
        <div id="chat">
            <div id="chat-header">{conversationTitle}</div>

            <div id="message-list">
                {loadingHistory && <div className="muted">Loading…</div>}
                {!loadingHistory && messages.length === 0 && !running && (
                    <div className="muted">No messages yet. Say hello.</div>
                )}
                {messages.map((m, i) => (
                    <div key={i} className={'message ' + m.role}>
                        <div className="message-role">{m.role === 'user' ? 'You' : 'Simon'}</div>
                        <div className="message-content">{m.content}</div>
                    </div>
                ))}
                {running && (
                    <div className="message assistant">
                        <div className="message-role">Simon</div>
                        <div className="message-content muted">Thinking…</div>
                    </div>
                )}
                <div ref={bottomRef}/>
            </div>

            {error && <div className="error-banner">{error}</div>}

            <div id="composer">
                <textarea
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder="Message Simon…"
                    rows={2}
                />
                {running ? (
                    <button id="cancel-button" onClick={handleCancel}>Cancel</button>
                ) : (
                    <button id="send-button" onClick={handleSend} disabled={!input.trim()}>Send</button>
                )}
            </div>
        </div>
    )
}

export default Chat
