/**
 * usePermission — returns a function to check if the current user has a given permission.
 * Permissions are stored in localStorage as part of the user object.
 *
 * Usage:
 *   const can = usePermission();
 *   can('create-produk') // true/false
 */
const usePermission = () => {
    const userData = localStorage.getItem('user');
    let mask = 0n;
    if (userData) {
        try {
            const parsed = JSON.parse(userData);
            if (parsed.permissions_mask) {
                mask = BigInt(parsed.permissions_mask);
            }
        } catch (e) {
            console.error("Failed to parse user data or mask", e);
        }
    }
    
    // requiredPermission is now a BigInt constant from PERMS
    return (requiredPermission) => {
        if (!requiredPermission) return true;
        return (mask & requiredPermission) !== 0n;
    };
};

export default usePermission;
