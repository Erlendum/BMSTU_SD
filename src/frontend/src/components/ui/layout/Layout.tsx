import { FC, PropsWithChildren } from 'react'
import styles from './Layout.module.scss'
import { Link } from 'react-router-dom'
import RegisterForm from '../form/RegisterForm'
import { UserService } from '../../../services/user.service'

const Layout: FC<PropsWithChildren<{ title?: string }>> = ({
	children,
	title
}) => {
	let isLogin = false
	let isAdmin = false
	isLogin = UserService.getCurrentLogin() != null

	if (isLogin) {
		return (
			<div className={styles.layout}>
				<header>
					<nav>
						{/*<h1>MUSIC STORE</h1>*/}
						<Link to='/'>Home</Link>
						<Link to='/discounts'>Discounts</Link>
						<Link to='/orders'>Orders</Link>
						<Link to='/comparison_list'>Comparison List</Link>
						<Link to='/logout'>Logout</Link>
					</nav>
					<h1>Music Store</h1>
				</header>
				{title && <h1 className={styles.heading}>{title}</h1>}
				{children}
			</div>
		)
	} else {
		return (
			<div className={styles.layout}>
				<header>
					<nav>
						{/*<h1>MUSIC STORE</h1>*/}
						<Link to='/'>Home</Link>
						<Link to='/discounts'>Discounts</Link>
						<Link to='/register'>Register</Link>
						<Link to='/login'>Login</Link>
					</nav>
				</header>
				{title && <h1 className={styles.heading}>{title}</h1>}
				{children}
			</div>
		)
	}
}

export default Layout
