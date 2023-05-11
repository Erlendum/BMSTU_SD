import React, { useEffect } from 'react'
import ReactDOM from 'react-dom/client'
import Home from './components/pages/home/Home'
import './index.css'
import { QueryClient, QueryClientProvider } from 'react-query'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import ComparisonList from './components/pages/comparisonList/ComparisonList'
import RegisterForm from './components/ui/form/RegisterForm'
import LoginForm from './components/ui/form/LoginForm'
import { UserService } from './services/user.service'
import Discounts from './components/pages/discounts/Discounts'

const queryClient = new QueryClient()

// @ts-ignore
const LogoutWrapper = ({ children }) => {
	useEffect(() => {
		UserService.logout()
	}, [])
	return children
}

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
	<QueryClientProvider client={queryClient}>
		<Router>
			<Routes>
				<Route path='/' element={<Home />} />
				<Route path='/discounts' element={<Discounts />} />
				<Route path='/comparison_list' element={<ComparisonList />} />
				<Route path='/register' element={<RegisterForm />} />
				<Route path='/login' element={<LoginForm />} />
				<Route
					path='/logout'
					element={<LogoutWrapper>Logout...</LogoutWrapper>}
				/>
			</Routes>
		</Router>
	</QueryClientProvider>
)
