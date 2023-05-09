import React from 'react'
import styles from './MessageBox.module.scss'

// @ts-ignore
const MessageBox = ({ type, children }) => {
	return (
		<div className={type == 'error' ? styles.error : styles.success}>
			{children}
		</div>
	)
}

export default MessageBox
