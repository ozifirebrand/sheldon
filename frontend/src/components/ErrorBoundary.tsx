import { Component, type ReactNode } from 'react'

interface Props {
  children: ReactNode
}

interface State {
  error: Error | null
}

/** Catches render errors and shows a fallback. */
export default class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null }

  static getDerivedStateFromError(error: Error): State {
    return { error }
  }

  render() {
    if (this.state.error) {
      return (
        <div style={{ padding: '2rem', fontFamily: 'system-ui', color: '#333' }}>
          <h1>Something went wrong</h1>
          <pre style={{ background: '#f0f0f0', padding: '1rem', overflow: 'auto' }}>
            {this.state.error.message}
          </pre>
          <p>Make sure the backend is running (see README).</p>
        </div>
      )
    }
    return this.props.children
  }
}
