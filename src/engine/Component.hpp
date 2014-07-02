#ifndef PLURALITY_ENGINE_COMPONENT_HPP
#define PLURALITY_ENGINE_COMPONENT_HPP

#include <string>
#include <map>

#include <boost/variant.hpp>
#include <boost/shared_ptr.hpp>

class Component {
	public:
		Component(const std::string& name) : mName(name) { }
		~Component() { }
		virtual void Start() { }
		inline void addValue(std::string name, const std::string& value);
		inline void addValue(std::string name, int value);

		virtual std::map<std::string, std::string> getPossibleValues() const { return {}; }
		const std::string& getName() const { return mName; }

	protected:
		std::string mName;
		std::map<std::string, boost::variant<std::string, int>> mValues;
};

typedef boost::shared_ptr<Component> ComponentPtr;

void Component::addValue(std::string name, const std::string& value)
{
	mValues[name] = value;
}

void Component::addValue(std::string name, int value)
{
	mValues[name] = value;
}


#endif

