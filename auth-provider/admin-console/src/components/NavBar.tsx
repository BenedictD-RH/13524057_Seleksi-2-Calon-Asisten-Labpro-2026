import './NavBar.css'

function getPathName(pageName: string): string {
    return '/' + pageName.toLowerCase()
}


function NavBar() {
    const path = window.location.pathname
    const pageList = ["Users", "Groups", "Apps"]
    return (
        <div className="navbar">
            {pageList.map((page, index) => (
                <li key={index}>
                    <a href={getPathName(page)}>
                        <button className={"navbutton" + (path == getPathName(page) ? " selected" : "")}>{page}</button>
                    </a>      
                </li>
            ))}
        </div>
    )
}

export default NavBar