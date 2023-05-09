import { FC, useEffect, useState } from 'react'
import Layout from '../../ui/layout/Layout'
import { useQuery } from 'react-query'
import { InstrumentService } from '../../../services/instrument.service'
import { UserService } from '../../../services/user.service'
import InstrumentItem from '../../ui/instrument-item/InstrumentItem'
import styles from '../cart/Cart.module.scss'
import { CartService } from '../../../services/cart.service'
import { OrderService } from '../../../services/order.service'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'

const Cart: FC = () => {
	const [error, setError] = useState('no error')
	const { data: cartInstruments, isLoading } = useQuery(
		['cartInstruments'],
		() => UserService.getCart(),
		{
			select: ({ cartInstruments }) => cartInstruments
		}
	)

	const { data: cart } = useQuery(['cart'], () => UserService.getCart(), {
		select: ({ cart }) => cart
	})

	const handleCheckout = async () => {
		let id
		await OrderService.create()
			.then(data => (id = data))
			.catch(error => {
				if (error.response) {
					setError(error.response.data.Error)
				}
			})
		if (error == 'no error') {
			toast.success('Order with id ' + id + ' was successfully created', {
				position: toast.POSITION.BOTTOM_LEFT
			})
		} else {
			toast.error('ERROR ' + error, {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
	}

	const displayCartInstruments = cartInstruments?.map(instrument => {
		return (
			<div>
				<InstrumentItem
					instrument={instrument}
					isCart={true}
					key={instrument.InstrumentId}
				/>
			</div>
		)
	})

	return (
		<Layout title='Cart'>
			<div className={styles.text}>Amount: {cart?.Amount}</div>
			<div className={styles.text}>Total Price: {cart?.TotalPrice}</div>
			<button
				className={styles.openBtn}
				onClick={handleCheckout}
				hidden={cartInstruments == null}
			>
				Checkout
			</button>
			<div className={styles.wrapper}>{displayCartInstruments}</div>
			<ToastContainer />
		</Layout>
	)
}

export default Cart
