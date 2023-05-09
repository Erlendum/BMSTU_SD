import React from 'react'
import styles from './Modal.module.scss'

// @ts-ignore
const Modal = ({ active, setActive, children }) => {
	return (
		<div
			className={active ? styles.modalActive : styles.modal}
			onClick={() => setActive(false)}
		>
			<div
				className={active ? styles.modal__contentActive : styles.modal__content}
				onClick={e => e.stopPropagation()}
			>
				{children}
			</div>
		</div>
	)
}

export default Modal
